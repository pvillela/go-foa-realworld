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
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticlesFeedSflT is the type of the stereotype instance for the service flow that
// queries for all articles authored by other users followed by
// the current user, with optional limit and offset pagination parameters.
type ArticlesFeedSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	in rpc.ArticlesFeedIn,
) (rpc.ArticlesOut, error)

// ArticlesFeedSflC is the function that constructs a stereotype instance of type
// ArticlesFeedSflT without hard-wired stereotype dependencies.
func ArticlesFeedSflC(
	cfgSrc DefaultSflCfgSrc,
	userGetByNameDaf daf.UserGetByNameDafT,
	articlesFeedDaf daf.ArticlesFeedDafT,
) ArticlesFeedSflT {
	db := cfgSrc()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		in rpc.ArticlesFeedIn,
	) (rpc.ArticlesOut, error) {
		username := reqCtx.Username
		zero := rpc.ArticlesOut{}

		user, err := userGetByNameDaf(ctx, tx, username)
		if err != nil {
			return zero, err
		}

		articlesPlus, err := articlesFeedDaf(ctx, tx, user.Id, in.Limit, in.Offset)
		if err != nil {
			return zero, err
		}

		articlesOut := rpc.ArticlesOut_FromModel(articlesPlus)

		return articlesOut, err
	})
}
