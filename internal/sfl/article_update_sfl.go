/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleUpdateSfl is the stereotype instance for the service flow that
// updates an article.
type ArticleUpdateSfl struct {
	ArticleGetAndCheckOwnerFl     fs.ArticleGetAndCheckOwnerFlT
	ArticleValidateBeforeUpdateBf fs.ArticleValidateBeforeUpdateBfT
	UserGetByNameDaf              fs.UserGetByNameDafT
	ArticleUpdateDaf              fs.ArticleUpdateDafT
	ArticleGetBySlugDaf           fs.ArticleGetBySlugDafT
	ArticleCreateDaf              fs.ArticleCreateDafT
}

// ArticleUpdateSflT is the function type instantiated by ArticleUpdateSfl.
type ArticleUpdateSflT = func(username, slug string, in rpc.ArticleUpdateIn) (rpc.ArticleOut, error)

func (s ArticleUpdateSfl) Make() ArticleUpdateSflT {
	return func(username string, slug string, in rpc.ArticleUpdateIn) (rpc.ArticleOut, error) {
		var zero rpc.ArticleOut

		article, rcArticle, err := s.ArticleGetAndCheckOwnerFl(slug, username)
		if err != nil {
			return zero, err
		}

		fieldsToUpdate := make(map[model.ArticleUpdatableField]interface{}, 3)
		if v := in.Article.Title; v != nil {
			fieldsToUpdate[model.Title] = *v
		}
		if v := in.Article.Description; v != nil {
			fieldsToUpdate[model.Description] = *v
		}
		if v := in.Article.Body; v != nil {
			fieldsToUpdate[model.Body] = *v
		}

		var newSlug string

		article, newSlug = article.Update(fieldsToUpdate)

		if err := s.ArticleValidateBeforeUpdateBf(article); err != nil {
			return zero, err
		}

		user, _, err := s.UserGetByNameDaf(username)
		if err != nil {
			return zero, err
		}

		// TODO: move some of this logic to a BF
		var savedArticle model.Article
		if newSlug == slug {
			savedArticle, _, err = s.ArticleUpdateDaf(article, rcArticle)
			if err != nil {
				return zero, err
			}
		} else {
			if _, _, err := s.ArticleGetBySlugDaf(newSlug); err == nil {
				return zero, fs.ErrDuplicateArticle
			}
			savedArticle, _, err = s.ArticleCreateDaf(article)
			if err != nil {
				return zero, err
			}
		}

		articleOut := rpc.ArticleOut{}.FromModel(user, savedArticle)
		return articleOut, err
	}
}
