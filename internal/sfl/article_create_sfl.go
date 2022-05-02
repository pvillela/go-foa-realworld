/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleCreateSflT is the type of the stereotype instance for the service flow that
// creates an article with hard-wired BF dependencies.
type ArticleCreateSflT = func(ctx context.Context, in rpc.ArticleCreateIn) (rpc.ArticleOut, error)

// ArticleCreateSflC is the function that constructs a stereotype instance of type
// ArticleCreateSflT with hard-wired BF dependencies.
func ArticleCreateSflC(
	beginTxn func(context string) db.Txn,
	userGetByNameDaf daf.UserGetByNameDafT,
	articleCreateDaf daf.ArticleCreateDafT,
	tagAddDaf daf.TagAddDafT,
) ArticleCreateSflT {
	articleValidateBeforeCreateBf := fs.ArticleValidateBeforeCreateBfI
	return ArticleCreateSflC0(
		beginTxn,
		userGetByNameDaf,
		articleCreateDaf,
		tagAddDaf,
		articleValidateBeforeCreateBf,
	)
}

// ArticleCreateSflC0 is the function that constructs a stereotype instance of type
// ArticleCreateSflT without hard-wired BF dependencies..
func ArticleCreateSflC0(
	beginTxn func(context string) db.Txn,
	userGetByNameDaf daf.UserGetByNameDafT,
	articleCreateDaf daf.ArticleCreateDafT,
	tagAddDaf daf.TagAddDafT,
	articleValidateBeforeCreateBf fs.ArticleValidateBeforeCreateBfT,
) ArticleCreateSflT {
	return func(ctx context.Context, in rpc.ArticleCreateIn) (rpc.ArticleOut, error) {
		username := web.ContextToRequestContext(ctx).Username
		txn := beginTxn("ArticleCreateSflS")
		defer txn.End()

		zero := rpc.ArticleOut{}

		user, _, err := userGetByNameDaf(username)
		if err != nil {
			return zero, err
		}

		article := in.ToArticle(user)

		if err := articleValidateBeforeCreateBf(article); err != nil {
			return zero, err
		}

		_, err = articleCreateDaf(article, txn)
		if err != nil {
			return zero, err
		}

		if err := tagAddDaf(article.TagList); err != nil {
			return zero, err
		}

		articleOut := rpc.ArticleOut_FromModel(user, article)
		return articleOut, err
	}
}
