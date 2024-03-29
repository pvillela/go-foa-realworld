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
	"github.com/pvillela/go-foa-realworld/experimental/fl"
	"github.com/pvillela/go-foa-realworld/experimental/rpc"
)

// ArticleGetSflT is the type of the stereotype instance for the service flow that
// retrieves an article.
type ArticleGetSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	slug string,
) (rpc.ArticleOut, error)

// ArticleGetSflC0 is the function that constructs a stereotype instance of type
// ArticleGetSflT without hard-wired stereotype dependencies.
func ArticleGetSflC0(
	cfgSrc DefaultSflCfgSrc,
	articleAndUserGetFl fl.ArticleAndUserGetFlT,
) ArticleGetSflT {
	db := cfgSrc()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		slug string,
	) (rpc.ArticleOut, error) {
		username := reqCtx.Username

		article, _, err := articleAndUserGetFl(ctx, tx, slug, username)
		if err != nil {
			return rpc.ArticleOut{}, err
		}

		articleOut := rpc.ArticleOut_FromModel(article)

		return articleOut, nil
	})
}
