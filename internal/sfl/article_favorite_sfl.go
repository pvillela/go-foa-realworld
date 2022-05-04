/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/fl"

	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleFavoriteSflT is the type of the stereotype instance for the service flow that
// designates an article as a favorite.
type ArticleFavoriteSflT = func(ctx context.Context, slug string) (rpc.ArticleOut, error)

// ArticleFavoriteSflC is the function that constructs a stereotype instance of type
// ArticleFavoriteSflT.
func ArticleFavoriteSflC(
	beginTxn func(context string) db.Txn,
	articleFavoriteFl fl.ArticleFavoriteFlT,
) ArticleFavoriteSflT {
	return func(ctx context.Context, slug string) (rpc.ArticleOut, error) {
		username := web.ContextToRequestContext(ctx).Username

		txn := beginTxn("ArticleCreateSflS")
		defer txn.End()

		var zero rpc.ArticleOut
		pwUser, pwArticle, err := articleFavoriteFl(username, slug, true, txn)
		if err != nil {
			return zero, err
		}
		articleOut := rpc.ArticleOut_FromModel(pwUser.Entity, pwArticle.Entity)
		return articleOut, err
	}
}
