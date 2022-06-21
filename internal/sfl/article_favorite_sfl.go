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

// ArticleFavoriteSflT is the type of the stereotype instance for the service flow that
// designates an article as a favorite.
type ArticleFavoriteSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	slug string,
) (rpc.ArticleOut, error)

// ArticleFavoriteSflC is the function that constructs a stereotype instance of type
// ArticleFavoriteSflT with hard-wired stereotype dependencies.
func ArticleFavoriteSflC(
	cfgPvdr DefaultSflCfgPvdr,
) ArticleFavoriteSflT {
	return ArticleFavoriteSflC0(
		cfgPvdr,
		fl.ArticleAndUserGetFl,
		daf.FavoriteCreateDaf,
		daf.ArticleUpdateDaf,
	)
}

// TODO: consider reimplementing with daf.ArticleAdjustFavoritesCountDaf.
// ArticleFavoriteSflC0 is the function that constructs a stereotype instance of type
// ArticleFavoriteSflT without hard-wired stereotype dependencies.
func ArticleFavoriteSflC0(
	cfgPvdr DefaultSflCfgPvdr,
	articleAndUserGetFl fl.ArticleAndUserGetFlT,
	favoriteCreateDaf daf.FavoriteCreateDafT,
	articleUpdateDaf daf.ArticleUpdateDafT,
) ArticleFavoriteSflT {
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

		err = favoriteCreateDaf(ctx, tx, articlePlus.Id, user.Id)
		if err != nil {
			return zero, err
		}

		article := articlePlus.ToArticle()
		article.AdjustFavoriteCount(1)

		err = articleUpdateDaf(ctx, tx, &article)
		if err != nil {
			return rpc.ArticleOut{}, err
		}

		articleOut := rpc.ArticleOut_FromModel(articlePlus)
		return articleOut, err
	})
}
