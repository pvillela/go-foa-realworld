/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
)

// ArticleFavoriteFlT is the type of the stereotype instance for the flow that
// designates an article as a favorite or not.
type ArticleFavoriteFlT = func(
	username, slug string,
	favorite bool,
	txn db.Txn,
) (PwUser, PwArticle, error)

// ArticleFavoriteFlC is the function that constructs a stereotype instance of type
// ArticleFavoriteFlT.
func ArticleFavoriteFlC(
	userGetByNameDaf UserGetByNameDafT,
	articleGetBySlugDaf ArticleGetBySlugDafT,
	articleUpdateDaf ArticleUpdateDafT,
) ArticleFavoriteFlT {
	return func(
		username, slug string,
		favorite bool,
		txn db.Txn,
	) (PwUser, PwArticle, error) {
		user, rcUser, err := userGetByNameDaf(username)
		if err != nil {
			return PwUser{}, PwArticle{}, err
		}

		article, rcArticle, err := articleGetBySlugDaf(slug)
		if err != nil {
			return PwUser{}, PwArticle{}, err
		}

		article = article.UpdateFavoritedBy(user, favorite)

		rcArticle, err = articleUpdateDaf(article, rcArticle, txn)
		if err != nil {
			return PwUser{}, PwArticle{}, err
		}

		return PwUser{rcUser, user}, PwArticle{rcArticle, article}, nil
	}
}
