/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	log "github.com/sirupsen/logrus"
	"strings"
)

func main() {
	defer util.PanicLog(log.Fatal)

	log.SetLevel(log.DebugLevel)
	//var arr []any
	//fmt.Println(arr[0])

	ctx := context.Background()

	connStr := "postgres://testuser:testpassword@localhost:9999/testdb?sslmode=disable"
	pool, err := pgxpool.Connect(ctx, connStr)
	util.PanicOnError(err)

	ctxDb := dbpgx.CtxPgx{pool}
	ctx, err = ctxDb.SetPool(ctx)
	util.PanicOnError(err)

	db := dbpgx.Db{pool}

	tx, err := db.BeginTx(ctx)
	util.PanicOnError(err)
	cleanupTables(ctx, tx, "users", "articles", "tags", "followings", "favorites", "article_tags")
	err = tx.Commit(ctx)
	util.PanicOnError(err)

	userDafsExample(ctx, ctxDb)
	articleDafsExample(ctx, db)
	tagDafsExample(ctx, db)
}

func cleanupTables(ctx context.Context, tx pgx.Tx, tables ...string) {
	tablesStr := strings.Join(tables, ", ")
	sql := fmt.Sprintf("TRUNCATE %v", tablesStr)
	_, err := tx.Exec(ctx, sql)
	util.PanicOnError(err)
}
