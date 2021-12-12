/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fs

import "github.com/pvillela/go-foa-realworld/internal/arch/db"

// ArticleFavoriteFlS is the stereotype instance for the flow that
// designates an article as a favorite or not.
type ArticleFavoriteFlS struct {
	UserGetByNameDaf    UserGetByNameDafT
	ArticleGetBySlugDaf ArticleGetBySlugDafT
	ArticleUpdateDaf    ArticleUpdateDafT
}

// ArticleFavoriteFlT is the function type instantiated by fs.ArticleFavoriteFlS.
type ArticleFavoriteFlT = func(username, slug string, favorite bool, txn db.Txn) (PwUser, PwArticle, error)

func (s ArticleFavoriteFlS) Make() ArticleFavoriteFlT {
	return func(username, slug string, favorite bool, txn db.Txn) (PwUser, PwArticle, error) {
		user, rcUser, err := s.UserGetByNameDaf(username)
		if err != nil {
			return PwUser{}, PwArticle{}, err
		}

		article, rcArticle, err := s.ArticleGetBySlugDaf(slug)
		if err != nil {
			return PwUser{}, PwArticle{}, err
		}

		article = article.UpdateFavoritedBy(user, favorite)

		rcArticle, err = s.ArticleUpdateDaf(article, rcArticle, txn)
		if err != nil {
			return PwUser{}, PwArticle{}, err
		}

		return PwUser{rcUser, user}, PwArticle{rcArticle, article}, nil
	}
}
