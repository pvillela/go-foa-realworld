/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticlesListSflS is the stereotype instance for the service flow that
// retrieve recent articles based on a set of query parameters.
type ArticlesListSflS struct {
	UserGetByNameDaf            fs.UserGetByNameDafT
	ArticleGetRecentFilteredDaf fs.ArticleGetRecentFilteredDafT
}

// ArticlesListSflT is the function type instantiated by ArticlesListSflS.
type ArticlesListSflT = func(username string, in rpc.ArticlesListIn) (rpc.ArticlesOut, error)

func (s ArticlesListSflS) Make() ArticlesListSflT {
	return func(username string, in rpc.ArticlesListIn) (rpc.ArticlesOut, error) {
		var zero rpc.ArticlesOut
		var user model.User
		var err error

		if username != "" {
			user, _, err = s.UserGetByNameDaf(username)
			if err != nil {
				return zero, err
			}
		}

		articles, err := s.ArticleGetRecentFilteredDaf(in)
		if err != nil {
			return zero, err
		}

		articlesOut := rpc.ArticlesOut_FromModel(user, articles)
		return articlesOut, err
	}
}
