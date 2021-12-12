/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
)

// ArticleFavoriteSflS is the stereotype instance for the service flow that
// deletes an article.
type ArticleDeleteSflS struct {
	BeginTxn                  func(context string) db.Txn
	ArticleGetAndCheckOwnerFl fs.ArticleGetAndCheckOwnerFlT
	ArticleDeleteDaf          fs.ArticleDeleteDafT
}

// ArticleDeleteSflT is the function type instantiated by ArticleDeleteSflS.
type ArticleDeleteSflT = func(username, slug string) error

func (s ArticleDeleteSflS) Make() ArticleDeleteSflT {
	return func(username string, slug string) error {
		txn := s.BeginTxn("ArticleCreateSflS")
		defer txn.End()

		_, _, err := s.ArticleGetAndCheckOwnerFl(username, slug)
		if err != nil {
			return err
		}

		return s.ArticleDeleteDaf(slug, txn)
	}
}
