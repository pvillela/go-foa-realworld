/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package testutil

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

func ArticlePlusesToArticles(aps []model.ArticlePlus) []model.Article {
	return util.SliceMap(aps, func(ap model.ArticlePlus) model.Article {
		return ap.ToArticle()
	})
}

func CleanupAllTables(db dbpgx.Db, ctx context.Context) {
	ctx, err := dbpgx.WithTransaction(db, ctx, func(ctx context.Context, tx pgx.Tx) (context.Context, error) {
		dbpgx.CleanupTables(ctx, tx, "users", "articles", "tags", "followings", "favorites",
			"article_tags", "comments")
		return ctx, nil
	})
	errx.PanicOnError(err)
}
