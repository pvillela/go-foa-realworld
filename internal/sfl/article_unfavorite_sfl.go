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
	db dbpgx.Db,
) ArticleUnfavoriteSflT {
	articleAndUserGetFl := fl.ArticleAndUserGetFlI
	favoriteDeleteDaf := daf.FavoriteDeleteDafI
	articleUpdateDaf := daf.ArticleUpdateDafI
	return ArticleUnfavoriteSflC0(
		db,
		articleAndUserGetFl,
		favoriteDeleteDaf,
		articleUpdateDaf,
	)
}

// ArticleUnfavoriteSflC0 is the function that constructs a stereotype instance of type
// ArticleUnfavoriteSflT without hard-wired stereotype dependencies.
func ArticleUnfavoriteSflC0(
	db dbpgx.Db,
	articleAndUserGetFl fl.ArticleAndUserGetFlT,
	favoriteDeleteDaf daf.FavoriteDeleteDafT,
	articleUpdateDaf daf.ArticleUpdateDafT,
) ArticleUnfavoriteSflT {
	return func(
		ctx context.Context,
		reqCtx web.RequestContext,
		slug string,
	) (rpc.ArticleOut, error) {
		return dbpgx.Db_WithTransaction(db, ctx, func(
			ctx context.Context,
			tx pgx.Tx,
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
			article.IncrementFavoriteCount()

			err = articleUpdateDaf(ctx, tx, &article)
			if err != nil {
				return rpc.ArticleOut{}, err
			}

			articleOut := rpc.ArticleOut_FromModel(articlePlus)
			return articleOut, err
		})
	}
}
