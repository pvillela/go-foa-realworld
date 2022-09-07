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
	"github.com/pvillela/go-foa-realworld/experimental/arch/types"
	"github.com/pvillela/go-foa-realworld/experimental/arch/web"
	"github.com/pvillela/go-foa-realworld/experimental/fl"
	"github.com/pvillela/go-foa-realworld/experimental/daf"
)

// ArticleDeleteSflT is the type of the stereotype instance for the service flow that
// deletes an article.
type ArticleDeleteSflT = func(ctx context.Context, reqCtx web.RequestContext, slug string) (types.Unit, error)

// ArticleDeleteSflC0 is the function that constructs a stereotype instance of type
// ArticleDeleteSflT without hard-wired stereotype dependencies.
func ArticleDeleteSflC0(
	cfgSrc DefaultSflCfgSrc,
	articleGetAndCheckOwnerFl fl.ArticleGetAndCheckOwnerFlT,
	articleDeleteDaf daf.ArticleDeleteDafT,
) ArticleDeleteSflT {
	db := cfgSrc()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		slug string,
	) (types.Unit, error) {
		username := reqCtx.Username

		_, _, err := articleGetAndCheckOwnerFl(ctx, tx, slug, username)
		if err != nil {
			return types.UnitV, err
		}

		// Record existence is guaranteed by above code.
		err = articleDeleteDaf(ctx, tx, slug)
		return types.UnitV, err
	})
}
