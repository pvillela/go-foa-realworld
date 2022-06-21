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
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticlesListSflT is the type of the stereotype instance for the service flow that
// retrieve recent articles based on a set of query parameters.
type ArticlesListSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	in model.ArticleCriteria,
) (rpc.ArticlesOut, error)

// ArticlesListSflC is the function that constructs a stereotype instance of type
// ArticlesListSflT with hard-wired stereotype dependencies.
func ArticlesListSflC(
	cfgPvdr DefaultSflCfgPvdr,
) ArticlesListSflT {
	return ArticlesListSflC0(
		cfgPvdr,
		daf.UserGetByNameExplicitTxDaf,
		daf.ArticlesListDaf,
	)
}

// ArticlesListSflC0 is the function that constructs a stereotype instance of type
// ArticlesListSflT without hard-wired stereotype dependencies.
func ArticlesListSflC0(
	cfgPvdr DefaultSflCfgPvdr,
	userGetByNameDaf daf.UserGetByNameExplicitTxDafT,
	articlesListDaf daf.ArticlesListDafT,
) ArticlesListSflT {
	db := cfgPvdr()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		in model.ArticleCriteria,
	) (rpc.ArticlesOut, error) {
		username := reqCtx.Username
		zero := rpc.ArticlesOut{}

		user, _, err := userGetByNameDaf(ctx, tx, username)
		if err != nil {
			return zero, err
		}

		articlesPlus, err := articlesListDaf(ctx, tx, user.Id, in)
		if err != nil {
			return zero, err
		}

		articlesOut := rpc.ArticlesOut_FromModel(articlesPlus)

		return articlesOut, err
	})
}
