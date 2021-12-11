/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleGetSfl is the stereotype instance for the service flow that
// retrieves an article.
type ArticleGetSfl struct {
	BeginTxn            func(context string) db.Txn
	UserGetByNameDaf    fs.UserGetByNameDafT
	ArticleGetBySlugDaf fs.ArticleGetBySlugDafT
}

// ArticleGetSflT is the function type instantiated by ArticleGetSfl.
type ArticleGetSflT = func(username string, slug string) (rpc.ArticleOut, error)

func (s ArticleGetSfl) Make() ArticleGetSflT {
	return func(username string, slug string) (rpc.ArticleOut, error) {
		txn := s.BeginTxn("ArticleCreateSfl")
		defer txn.End()

		var zero rpc.ArticleOut
		var user model.User
		var err error

		if username != "" {
			user, _, err = s.UserGetByNameDaf(username)
			if err != nil {
				return zero, err
			}
		}

		article, _, err := s.ArticleGetBySlugDaf(slug)
		if err != nil {
			return zero, err
		}

		articleOut := rpc.ArticleOut_FromModel(user, article)

		return articleOut, err
	}
}
