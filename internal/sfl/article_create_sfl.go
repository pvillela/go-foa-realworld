/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleCreateSflT is the type of the stereotype instance for the service flow that
// creates an article.
type ArticleCreateSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	in rpc.ArticleCreateIn,
) (rpc.ArticleOut, error)

// ArticleCreateSflC is the function that constructs a stereotype instance of type
// ArticleCreateSflT with hard-wired stereotype dependencies.
func ArticleCreateSflC(
	db dbpgx.Db,
	userGetByNameDaf daf.UserGetByNameDafT,
	articleCreateDaf daf.ArticleCreateDafT,
	tagAddDaf daf.TagCreateDafT,
) ArticleCreateSflT {
	articleValidateBeforeCreateBf := bf.ArticleValidateBeforeCreateBfI
	return ArticleCreateSflC0(
		beginTxn,
		userGetByNameDaf,
		articleCreateDaf,
		tagAddDaf,
		articleValidateBeforeCreateBf,
	)
}

// ArticleCreateSflC0 is the function that constructs a stereotype instance of type
// ArticleCreateSflT without hard-wired stereotype dependencies..
func ArticleCreateSflC0(
	db dbpgx.Db,
	userGetByNameDaf daf.UserGetByNameExplicitTxDafT,
	articleCreateDaf daf.ArticleCreateDafT,
	tagAddDaf daf.TagCreateDafT,
) ArticleCreateSflT {
	return func(
		ctx context.Context,
		reqCtx web.RequestContext,
		in rpc.ArticleCreateIn,
	) (rpc.ArticleOut, error) {
		return dbpgx.Db_WithTransaction(db, ctx, func(
			ctx context.Context,
			tx pgx.Tx,
		) (rpc.ArticleOut, error) {
			zero := rpc.ArticleOut{}
			username := reqCtx.Username

			user, _, err := userGetByNameDaf(ctx, tx, username)
			if err != nil {
				return rpc.ArticleOut{}, err
			}

			article := in.ToArticle(&user)
			err = articleCreateDaf(ctx, tx, &article)
			if err != nil {
				return rpc.ArticleOut{}, err
			}

		})
	}

	func(ctx context.Context, in rpc.ArticleCreateIn) (rpc.ArticleOut, error) {
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
