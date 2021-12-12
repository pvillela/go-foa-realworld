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

// ArticlesFeedSflS is the stereotype instance for the service flow that
// queries for all articles authored by other users followed by
// the current user, with optional limit and offset pagination parameters.
type ArticlesFeedSflS struct {
	UserGetByNameDaf              fs.UserGetByNameDafT
	ArticleGetRecentForAuthorsDaf fs.ArticleGetRecentForAuthorsDafT
}

// ArticlesFeedSflT is the function type instantiated by ArticlesFeedSflS.
type ArticlesFeedSflT = func(username string, in rpc.ArticlesFeedIn) (rpc.ArticlesOut, error)

func (s ArticlesFeedSflS) Make() ArticlesFeedSflT {
	return func(username string, in rpc.ArticlesFeedIn) (rpc.ArticlesOut, error) {
		var zero rpc.ArticlesOut
		var user model.User
		var err error

		if username == "" {
			return zero, fs.ErrNotAuthenticated.Make(nil)
		}

		user, _, err = s.UserGetByNameDaf(username)
		if err != nil {
			return zero, err
		}

		articles, err := s.ArticleGetRecentForAuthorsDaf(user.FollowedNames, in.Limit, in.Offset)
		if err != nil {
			return zero, err
		}

		articlesOut := rpc.ArticlesOut_FromModel(user, articles)
		return articlesOut, err
	}
}
