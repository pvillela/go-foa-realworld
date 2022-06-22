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
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx/dbpgxtest"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

func ArticlePlusesToArticles(aps []model.ArticlePlus) []model.Article {
	return util.SliceMap(aps, func(ap model.ArticlePlus) model.Article {
		return ap.ToArticle()
	})
}

func CleanupAllTables(db dbpgx.Db, ctx context.Context) {
	_, err := dbpgx.WithTransaction(db, ctx, func(ctx context.Context, tx pgx.Tx) (context.Context, error) {
		dbpgxtest.CleanupTables(ctx, tx, "users", "articles", "tags", "followings", "favorites",
			"article_tags", "comments")
		return ctx, nil
	})
	errx.PanicOnError(err)
}

func UserGetByName(db dbpgx.Db, ctx context.Context, username string) (model.User, error) {
	f := func(ctx context.Context, tx pgx.Tx) (model.User, error) {
		user, _, err := daf.UserGetByNameExplicitTxDaf(ctx, tx, username)
		return user, err
	}
	return dbpgx.WithTransaction(db, ctx, f)
}

func ArticleGetBySlug(db dbpgx.Db, ctx context.Context, currUsername string, slug string) (model.Article, error) {
	user, err := UserGetByName(db, ctx, currUsername)
	if err != nil {
		return model.Article{}, err
	}

	f := func(ctx context.Context, tx pgx.Tx) (model.Article, error) {
		articlePlus, err := daf.ArticleGetBySlugDaf(ctx, tx, user.Id, slug)
		if err != nil {
			return model.Article{}, err
		}
		article := articlePlus.ToArticle()
		return article, err
	}
	return dbpgx.WithTransaction(db, ctx, f)
}
