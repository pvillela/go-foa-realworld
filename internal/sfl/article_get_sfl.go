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

// ArticleGetSflT is the type of the stereotype instance for the service flow that
// retrieves an article.
type ArticleGetSflT = func(username string, slug string) (rpc.ArticleOut, error)

// ArticleGetSflC is the function that constructs a stereotype instance of type
// ArticleGetSflT.
func ArticleGetSflC(
	beginTxn func(context string) db.Txn,
	userGetByNameDaf fs.UserGetByNameDafT,
	articleGetBySlugDaf fs.ArticleGetBySlugDafT,
) ArticleGetSflT {
	return func(username string, slug string) (rpc.ArticleOut, error) {
		txn := beginTxn("ArticleCreateSflS")
		defer txn.End()

		var zero rpc.ArticleOut
		var user model.User
		var err error

		if username != "" {
			user, _, err = userGetByNameDaf(username)
			if err != nil {
				return zero, err
			}
		}

		article, _, err := articleGetBySlugDaf(slug)
		if err != nil {
			return zero, err
		}

		articleOut := rpc.ArticleOut_FromModel(user, article)

		return articleOut, err
	}
}
