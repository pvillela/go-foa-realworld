/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/experimental/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/experimental/arch/web"
	"github.com/pvillela/go-foa-realworld/experimental/daf"
	rpc2 "github.com/pvillela/go-foa-realworld/experimental/rpc"
)

// ArticlesListSflT is the type of the stereotype instance for the service flow that
// retrieve recent articles based on a set of query parameters.
type ArticlesListSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	in rpc2.ArticleCriteria,
) (rpc2.ArticlesOut, error)

// ArticlesListSflC is the function that constructs a stereotype instance of type
// ArticlesListSflT with hard-wired stereotype dependencies.
func ArticlesListSflC(
	cfgSrc DefaultSflCfgSrc,
) ArticlesListSflT {
	return ArticlesListSflC0(
		cfgSrc,
		daf.UserGetByNameExplicitTxDaf,
		daf.ArticlesListDaf,
	)
}

// ArticlesListSflC0 is the function that constructs a stereotype instance of type
// ArticlesListSflT without hard-wired stereotype dependencies.
func ArticlesListSflC0(
	cfgSrc DefaultSflCfgSrc,
	userGetByNameDaf daf.UserGetByNameExplicitTxDafT,
	articlesListDaf daf.ArticlesListDafT,
) ArticlesListSflT {
	db := cfgSrc()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		in rpc2.ArticleCriteria,
	) (rpc2.ArticlesOut, error) {
		username := reqCtx.Username
		zero := rpc2.ArticlesOut{}

		user, _, err := userGetByNameDaf(ctx, tx, username)
		if err != nil {
			return zero, err
		}

		articlesPlus, err := articlesListDaf(ctx, tx, user.Id, in)
		if err != nil {
			return zero, err
		}

		articlesOut := rpc2.ArticlesOut_FromModel(articlesPlus)

		return articlesOut, err
	})
}
