/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"
)

var connStr = "postgres://testuser:testpassword@localhost:9999/testdb?sslmode=disable"

type TestPair struct {
	Name string
	Func func(t *testing.T, db dbpgx.Db, ctx context.Context)
}

func dafTester0(testPairs []TestPair) func(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, connStr)
	errx.PanicOnError(err)

	ctxDb := dbpgx.CtxPgx{pool}
	ctx, err = ctxDb.SetPool(ctx)
	errx.PanicOnError(err)

	db := dbpgx.Db{pool}

	tx, err := db.BeginTx(ctx)
	errx.PanicOnError(err)
	cleanupTables(ctx, tx, "users", "articles", "tags", "followings", "favorites",
		"article_tags", "comments")
	setupData(ctx, tx)
	err = tx.Commit(ctx)
	errx.PanicOnError(err)

	return func(t *testing.T) {
		defer errx.PanicLog(log.Fatal)
		defer pool.Close()

		for _, p := range testPairs {
			testFunc := func(t *testing.T) {
				p.Func(t, db, ctx)
			}
			t.Run(p.Name, testFunc)
		}
	}
}

func dafTester(t *testing.T, testPairs []TestPair) {
	defer errx.PanicLog(log.Fatal)

	log.SetLevel(log.DebugLevel)

	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, connStr)
	errx.PanicOnError(err)

	ctxDb := dbpgx.CtxPgx{pool}
	ctx, err = ctxDb.SetPool(ctx)
	errx.PanicOnError(err)

	db := dbpgx.Db{pool}
	defer pool.Close()

	tx, err := db.BeginTx(ctx)
	errx.PanicOnError(err)
	cleanupTables(ctx, tx, "users", "articles", "tags", "followings", "favorites",
		"article_tags", "comments")
	setupData(ctx, tx)
	err = tx.Commit(ctx)
	errx.PanicOnError(err)

	for _, p := range testPairs {
		testFunc := func(t *testing.T) {
			p.Func(t, db, ctx)
		}
		t.Run(p.Name, testFunc)
	}
}

func cleanupTables(ctx context.Context, tx pgx.Tx, tables ...string) {
	tablesStr := strings.Join(tables, ", ")
	sql := fmt.Sprintf("TRUNCATE %v", tablesStr)
	_, err := tx.Exec(ctx, sql)
	errx.PanicOnError(err)
}

func setupData(ctx context.Context, tx pgx.Tx) {
	for i, _ := range users {
		recCtx, err := daf.UserCreateExplicitTxDafI(ctx, tx, &users[i])
		errx.PanicOnError(err)
		_, _ = spew.Printf("user from Create: %v\n", users[i])
		fmt.Println("recCtx from Create:", recCtx)
	}
	for i, _ := range articles {
		articles[i].AuthorId = users[1].Id
		err := daf.ArticleCreateDafI(ctx, tx, &articles[i])
		errx.PanicOnError(err)
		fmt.Println("article from Create:", articles[i])
	}
}

//func SubtestArticleCreateDafI(t *testing.T, db dbpgx.Db, ctx context.Context) {
//	_, err := dbpgx.Db_WithTransaction(db, ctx, func(ctx context.Context, tx pgx.Tx) (arch.Unit, error) {
//		currUser := users[0]
//		authors := []model.User{users[1], users[1]}
//		author := authors[0]
//
//		{
//			criteria := model.ArticleCriteria{
//				Tag:         nil,
//				Author:      &author.Username,
//				FavoritedBy: nil,
//				Limit:       nil,
//				Offset:      nil,
//			}
//			articlePluses, err := daf.ArticlesListDafI(ctx, tx, currUser.Id, criteria)
//			errx.PanicOnError(err)
//
//			returned := util.SliceMap(articlePluses, articlePlusToCore)
//			expected := util.SliceMap(articles, articleToCore(authors, false))
//
//			//fmt.Println("\ncoreInfoReturned:", coreInfoReturned)
//			//fmt.Println("\ncoreInfoExpected:", coreInfoExpected)
//
//			assert.ElementsMatch(t, expected, returned)
//		}
//
//		{
//			criteria := model.ArticleCriteria{
//				Tag:         util.PointerFromValue("FOOTAG"),
//				Author:      nil,
//				FavoritedBy: nil,
//				Limit:       nil,
//				Offset:      nil,
//			}
//			articlePluses, err := daf.ArticlesListDafI(ctx, tx, currUser.Id, criteria)
//			errx.PanicOnError(err)
//
//			returned := util.SliceMap(articlePluses, articlePlusToCore)
//			var expected []model.ArticlePlus
//
//			//fmt.Println("\ncore info returned:", returned)
//			//fmt.Println("\ncore info expected:", expected)
//
//			assert.ElementsMatch(t, expected, returned)
//		}
//
//		{
//			articleFromDb, err := daf.ArticleGetBySlugDafI(ctx, tx, currUser.Id, articles[1].Slug)
//			errx.PanicOnError(err)
//
//			returned := articlePlusToCore(-1, articleFromDb)
//			expected := articleToCore(authors, false)(0, articles[1])
//
//			//_, _ = spew.Println("\nArticleGetBySlugDaf:", articleFromDb)
//
//			assert.Equal(t, expected, returned)
//		}
//
//		{
//			pArticle := &articles[0]
//			pArticle.Title = "A very interesting subject"
//			err := daf.ArticleUpdateDafI(ctx, tx, pArticle)
//			errx.PanicOnError(err)
//
//			articleFromDb, err := daf.ArticleGetBySlugDafI(ctx, tx, currUser.Id, pArticle.Slug)
//			errx.PanicOnError(err)
//
//			returned := articlePlusToCore(-1, articleFromDb)
//			expected := articleToCore(authors, false)(0, *pArticle)
//
//			//_, _ = spew.Println("\nAfter update:", articleFromDb)
//
//			assert.Equal(t, expected, returned)
//		}
//
//		{
//			articlePluses, err := daf.ArticlesFeedDafI(ctx, tx, currUser.Id, nil, nil)
//			errx.PanicOnError(err)
//
//			returned := util.SliceMap(articlePluses, articlePlusToCore)
//			var expected []model.ArticlePlus
//
//			//_, _ = spew.Println("\nArticlesFeedDaf returned:", returned)
//			//_, _ = spew.Println("\nArticlesFeedDaf expected:", expected)
//
//			assert.ElementsMatch(t, expected, returned)
//		}
//
//		return arch.Void, nil
//	})
//	errx.PanicOnError(err)
//}
