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
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"
)

var connStr = "postgres://testuser:testpassword@localhost:9999/testdb?sslmode=disable"

// dafTester runs txnlSubtest by providing a database with repeatably initialized data.
func dafTester(t *testing.T, txnlSubtest dbpgx.TransactionalSubtest) {
	log.SetLevel(log.InfoLevel)

	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, connStr)
	errx.PanicOnError(err)

	db := dbpgx.Db{pool}
	defer pool.Close()

	ctx, err = dbpgx.WithTransaction(db, ctx, func(ctx context.Context, tx pgx.Tx) (context.Context, error) {
		cleanupTables(ctx, tx, "users", "articles", "tags", "followings", "favorites",
			"article_tags", "comments")
		setupData(ctx, tx)
		return ctx, nil
	})

	txnlSubtest(db, ctx, t)
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
		recCtxUsers[i] = recCtx
		_, _ = spew.Printf("user from Create: %v\n", users[i])
		log.Debug("recCtx from Create:", recCtx)
	}

	authors := []model.User{users[1], users[1]}
	if len(authors) != len(articles) {
		panic("len(authors) != len(articles)")
	}

	for i, _ := range articles {
		articles[i].AuthorId = authors[i].Id
		err := daf.ArticleCreateDafI(ctx, tx, &articles[i])
		errx.PanicOnError(err)
		log.Debug("article from Create:", articles[i])
	}
}
