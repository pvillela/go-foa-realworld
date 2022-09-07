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
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/config"
	"github.com/pvillela/go-foa-realworld/internal/fl"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
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
	db := cfgSrc.Get()
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

///////////////////
// Config logic

var ArticleDeleteSflCfgSrc = config.MakeConfigSource[DefaultSflCfgInfo](nil)

func articleDeleteSflCfgAdapter(appCfg config.AppCfgInfo) DefaultSflCfgSrc {
	return util.Todo[DefaultSflCfgSrc]()
}

// ArticleDeleteSflC is the function that constructs a stereotype instance of type
// ArticleDeleteSflT with hard-wired stereotype dependencies.
func ArticleDeleteSflC() ArticleDeleteSflT {
	return ArticleDeleteSflC0(
		ArticleDeleteSflCfgSrc,
		fl.ArticleGetAndCheckOwnerFl,
		daf.ArticleDeleteDaf,
	)
}
