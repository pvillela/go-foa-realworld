/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleCreateSflS is the stereotype instance for the service flow that
// creates an article.
type ArticleCreateSflS struct {
	BeginTxn         func(context string) db.Txn
	UserGetByNameDaf fs.UserGetByNameDafT
	ArticleCreateDaf fs.ArticleCreateDafT
	TagAddDaf        fs.TagAddDafT
}

// ArticleCreateSflT is the function type instantiated by ArticleCreateSflS.
type ArticleCreateSflT = func(username string, in rpc.ArticleCreateIn) (rpc.ArticleOut, error)

func (s ArticleCreateSflS) Make() ArticleCreateSflT {
	articleValidateBeforeCreateBf := fs.ArticleValidateBeforeCreateBfI
	return func(username string, in rpc.ArticleCreateIn) (rpc.ArticleOut, error) {
		txn := s.BeginTxn("ArticleCreateSflS")
		defer txn.End()

		zero := rpc.ArticleOut{}

		user, _, err := s.UserGetByNameDaf(username)
		if err != nil {
			return zero, err
		}

		article := in.ToArticle(user)

		if err := articleValidateBeforeCreateBf(article); err != nil {
			return zero, err
		}

		_, err = s.ArticleCreateDaf(article, txn)
		if err != nil {
			return zero, err
		}

		if err := s.TagAddDaf(article.TagList); err != nil {
			return zero, err
		}

		articleOut := rpc.ArticleOut_FromModel(user, article)
		return articleOut, err
	}
}
