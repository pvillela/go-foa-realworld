/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/newdaf"

	"github.com/pvillela/go-foa-realworld/internal/arch/web"

	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleUpdateSflT is the type of the stereotype instance for the service flow that
// updates an article, with hard-wired BF dependencies.
type ArticleUpdateSflT = func(context.Context, rpc.ArticleUpdateIn) (rpc.ArticleOut, error)

// ArticleUpdateSflC is the function that constructs a stereotype instance of type
// ArticleUpdateSflT with hard-wired BF dependencies.
func ArticleUpdateSflC(
	beginTxn func(context string) db.Txn,
	articleGetAndCheckOwnerFl fs.ArticleGetAndCheckOwnerFlT,
	userGetByNameDaf newdaf.UserGetByNameDafT,
	articleUpdateDaf newdaf.ArticleUpdateDafT,
	articleGetBySlugDaf newdaf.ArticleGetBySlugDafT,
	articleCreateDaf newdaf.ArticleCreateDafT,
	articleDeleteDaf newdaf.ArticleDeleteDafT,
) ArticleUpdateSflT {
	articleValidateBeforeUpdateBf := fs.ArticleValidateBeforeUpdateBfI
	return ArticleUpdateSflC0(
		beginTxn,
		articleGetAndCheckOwnerFl,
		userGetByNameDaf,
		articleUpdateDaf,
		articleGetBySlugDaf,
		articleCreateDaf,
		articleDeleteDaf,
		articleValidateBeforeUpdateBf,
	)
}

// ArticleUpdateSflC0 is the function that constructs a stereotype instance of type
// ArticleUpdateSflT without hard-wired BF dependencies.
func ArticleUpdateSflC0(
	beginTxn func(context string) db.Txn,
	articleGetAndCheckOwnerFl fs.ArticleGetAndCheckOwnerFlT,
	userGetByNameDaf newdaf.UserGetByNameDafT,
	articleUpdateDaf newdaf.ArticleUpdateDafT,
	articleGetBySlugDaf newdaf.ArticleGetBySlugDafT,
	articleCreateDaf newdaf.ArticleCreateDafT,
	articleDeleteDaf newdaf.ArticleDeleteDafT,
	articleValidateBeforeUpdateBf fs.ArticleValidateBeforeUpdateBfT,
) ArticleUpdateSflT {
	return func(ctx context.Context, in rpc.ArticleUpdateIn) (rpc.ArticleOut, error) {
		username := web.ContextToRequestContext(ctx).Username

		slug := in.Article.Slug
		txn := beginTxn("ArticleCreateSflS")
		defer txn.End()

		var zero rpc.ArticleOut

		article, rc, err := articleGetAndCheckOwnerFl(slug, username)
		if err != nil {
			return zero, err
		}

		updateSrc := model.ArticlePatch{
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
