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
	"github.com/pvillela/go-foa-realworld/internal/fl"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleUnfavoriteSflT is the type of the stereotype instance for the service flow that
// designates an article as no longer a favorite.
type ArticleUnfavoriteSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	slug string,
) (rpc.ArticleOut, error)

// ArticleUnfavoriteSflC is the function that constructs a stereotype instance of type
// ArticleUnfavoriteSflT with hard-wired stereotype dependencies.
func ArticleUnfavoriteSflC(
	cfgPvdr DefaultSflCfgPvdr,
) ArticleUnfavoriteSflT {
	return ArticleUnfavoriteSflC0(
		cfgPvdr,
		fl.ArticleAndUserGetFl,
		daf.FavoriteDeleteDaf,
		daf.ArticleUpdateDaf,
	)
}

// ArticleUnfavoriteSflC0 is the function that constructs a stereotype instance of type
// ArticleUnfavoriteSflT without hard-wired stereotype dependencies.
func ArticleUnfavoriteSflC0(
	cfgPvdr DefaultSflCfgPvdr,
	articleAndUserGetFl fl.ArticleAndUserGetFlT,
	favoriteDeleteDaf daf.FavoriteDeleteDafT,
	articleUpdateDaf daf.ArticleUpdateDafT,
) ArticleUnfavoriteSflT {
	db := cfgPvdr()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		slug string,
	) (rpc.ArticleOut, error) {
		username := reqCtx.Username
		var zero rpc.ArticleOut

		articlePlus, user, err := articleAndUserGetFl(ctx, tx, slug, username)
		if err != nil {
			return zero, err
		}

		err = favoriteDeleteDaf(ctx, tx, articlePlus.Id, user.Id)
		if err != nil {
			return zero, err
		}

		article := articlePlus.ToArticle()
		article.AdjustFavoriteCount(-1)

		err = articleUpdateDaf(ctx, tx, &article)
		if err != nil {
			return rpc.ArticleOut{}, err
		}

		articleOut := rpc.ArticleOut_FromModel(articlePlus)
		return articleOut, err
	})
}
