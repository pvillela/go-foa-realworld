/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleCreateSfl is the stereotype instance for the service flow that
// creates an article.
type ArticleCreateSfl struct {
	UserGetByNameDaf              fs.UserGetByNameDafT
	ArticleValidateBeforeCreateBf fs.ArticleValidateBeforeCreateBfT
	ArticleCreateDaf              fs.ArticleCreateDafT
	TagAddDaf                     fs.TagAddDafT
}

// ArticleCreateSflT is the function type instantiated by ArticleCreateSfl.
type ArticleCreateSflT = func(username string, in rpc.ArticleCreateIn) (rpc.ArticleOut, error)

func (s ArticleCreateSfl) Make() ArticleCreateSflT {
	return func(username string, in rpc.ArticleCreateIn) (rpc.ArticleOut, error) {
		zero := rpc.ArticleOut{}

		user, _, err := s.UserGetByNameDaf(username)
		if err != nil {
			return zero, err
		}

		article := in.ToArticle(user)

		if err := s.ArticleValidateBeforeCreateBf(article); err != nil {
			return zero, err
		}

		_, err = s.ArticleCreateDaf(article)
		if err != nil {
			return zero, err
		}

		if err := s.TagAddDaf(article.TagList); err != nil {
			return zero, err
		}

		articleOut := rpc.ArticleOut{}.FromModel(user, article)
		return articleOut, err
	}
}
