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

// ArticleUpdateSflT is the type of the stereotype instance for the service flow that
// updates an article.
type ArticleUpdateSflT = func(username, slug string, in rpc.ArticleUpdateIn) (rpc.ArticleOut, error)

// ArticleUpdateSflTC is the function that constructs a stereotype instance of type
// ArticleUpdateSflT.
func ArticleUpdateSflTC(
	beginTxn func(context string) db.Txn,
	articleGetAndCheckOwnerFl fs.ArticleGetAndCheckOwnerFlT,
	userGetByNameDaf fs.UserGetByNameDafT,
	articleUpdateDaf fs.ArticleUpdateDafT,
	articleGetBySlugDaf fs.ArticleGetBySlugDafT,
	articleCreateDaf fs.ArticleCreateDafT,
	articleDeleteDaf fs.ArticleDeleteDafT,
) ArticleUpdateSflT {
	articleValidateBeforeUpdateBf := fs.ArticleValidateBeforeUpdateBfI
	return func(username string, slug string, in rpc.ArticleUpdateIn) (rpc.ArticleOut, error) {
		txn := beginTxn("ArticleCreateSflS")
		defer txn.End()

		var zero rpc.ArticleOut

		article, rc, err := articleGetAndCheckOwnerFl(slug, username)
		if err != nil {
			return zero, err
		}

		updateSrc := model.ArticleUpdateSrc{
			Title:       in.Article.Title,
			Description: in.Article.Description,
			Body:        in.Article.Body,
		}

		article = article.Update(updateSrc)
		newSlug := article.Slug

		if err := articleValidateBeforeUpdateBf(article); err != nil {
			return zero, err
		}

		user, _, err := userGetByNameDaf(username)
		if err != nil {
			return zero, err
		}

		if newSlug == slug {
			_, err = articleUpdateDaf(article, rc, txn)
			if err != nil {
				return zero, err
			}
		} else {
			_, err = articleCreateDaf(article, txn)
			if err != nil {
				return zero, err
			}
			err = articleDeleteDaf(slug, txn)
			if err != nil {
				return zero, err
			}
		}

		articleOut := rpc.ArticleOut_FromModel(user, article)
		return articleOut, err
	}
}
