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
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	log "github.com/sirupsen/logrus"
	"strings"
)

func main() {
	defer errx.PanicLog(log.Fatal)

	log.SetLevel(log.DebugLevel)
	//var arr []any
	//fmt.Println(arr[0])

	ctx := context.Background()

	connStr := "postgres://testuser:testpassword@localhost:9999/testdb?sslmode=disable"
	pool, err := pgxpool.Connect(ctx, connStr)
	errx.PanicOnError(err)

	db := dbpgx.Db{pool}

	tx, err := db.BeginTx(ctx)
	errx.PanicOnError(err)
	cleanupTables(ctx, tx, "users", "articles", "tags", "followings", "favorites",
		"article_tags", "comments")
	err = tx.Commit(ctx)
	errx.PanicOnError(err)

	userDafsExample(ctx, db)
	articleDafsExample(ctx, db)
	tagDafsExample(ctx, db)
	followingDafsExample(ctx, db)
	favoriteDafsExample(ctx, db)
	commentDafsExample(ctx, db)
}

func cleanupTables(ctx context.Context, tx pgx.Tx, tables ...string) {
	tablesStr := strings.Join(tables, ", ")
	sql := fmt.Sprintf("TRUNCATE %v", tablesStr)
	_, err := tx.Exec(ctx, sql)
	errx.PanicOnError(err)
}
