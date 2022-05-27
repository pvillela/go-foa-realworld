/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf/daftest"
	log "github.com/sirupsen/logrus"
)

var connStr = "postgres://testuser:testpassword@localhost:9999/testdb?sslmode=disable"

func main() {
	defer errx.PanicLog(log.Fatal)

	log.SetLevel(log.DebugLevel)
	//var arr []any
	//fmt.Println(arr[0])

	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, connStr)
	errx.PanicOnError(err)

	//ctxDb := dbpgx.CtxPgx{pool}
	//ctx, err = ctxDb.SetPool(ctx)
	//errx.PanicOnError(err)

	db := dbpgx.Db{pool}

	tx, err := db.BeginTx(ctx)
	errx.PanicOnError(err)
	daftest.CleanupTables(ctx, tx, "users", "articles", "tags", "followings", "favorites",
		"article_tags", "comments")
	daftest.SetupData(ctx, tx)
	err = tx.Commit(ctx)
	errx.PanicOnError(err)

	daftest.UserDafsSubt(db, ctx, nil)
	//daftest.UserDafsSubt1(db, ctx)

}
