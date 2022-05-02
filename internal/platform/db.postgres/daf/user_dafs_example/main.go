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
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer util.PanicLog(log.Fatal)
	myBio := "I am me."
	user := model.User{
		Username:     "pvillela",
		Email:        "foo@bar.com",
		PasswordHash: "dakfljads0fj",
		Bio:          &myBio,
		ImageLink:    "",
	}

	ctx := context.Background()

	connStr := "postgres://testuser:testpassword@localhost:9999/testdb?sslmode=disable"
	pool, err := pgxpool.Connect(ctx, connStr)
	util.PanicOnError(err)
	ctxDb := dbpgx.CtxPgx{pool}
	ctx, err = ctxDb.SetPool(ctx)
	util.PanicOnError(err)

	ctx, err = ctxDb.BeginTx(ctx)
	util.PanicOnError(err)
	//fmt.Println("ctx", ctx)

	tx, err := dbpgx.GetCtxTx(ctx)
	util.PanicOnError(err)
	cleanupTable(ctx, tx, "users")

	recCtx, err := daf.UserCreateDaf(ctx, &user)
	util.PanicOnError(err)
	fmt.Println("recCtx from Create:", recCtx)

	userFromDb, recCtx, err := daf.UserGetByNameDaf(ctx, "pvillela")
	util.PanicOnError(err)
	fmt.Println("\nUserGetByNameDaf:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx)

	userFromDb, recCtx, err = daf.UserGetByEmailDaf(ctx, "foo@bar.com")
	util.PanicOnError(err)
	fmt.Println("\nUserGetByEmailDaf:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx)

	readManySql := "SELECT * FROM users"
	pwUsers, err := dbpgx.ReadMany[daf.PwUser](ctx, tx, readManySql)
	fmt.Println("pwUsers:", pwUsers)
	pwUserPtrs, err := dbpgx.ReadMany[*daf.PwUser](ctx, tx, readManySql)
	fmt.Println("pwUserPtrs:", pwUserPtrs)
	fmt.Println("*pwUserPtrs[0]:", *pwUserPtrs[0])

	ctx, err = ctxDb.Commit(ctx)
	util.PanicOnError(err)

	ctx, err = ctxDb.BeginTx(ctx)
	util.PanicOnError(err)

	user.ImageLink = "https://xyz.com"
	recCtx, err = daf.UserUpdateDaf(ctx, user, recCtx)
	util.PanicOnError(err)
	fmt.Println("\nUserUpdateDaf:", user)
	fmt.Println("recCtx from Update:", recCtx)

	_, err = ctxDb.Commit(ctx)
	util.PanicOnError(err)
}

func cleanupTable(ctx context.Context, tx pgx.Tx, table string) {
	sql := fmt.Sprintf("TRUNCATE %v", table)
	_, err := tx.Exec(ctx, sql)
	util.PanicOnError(err)
}
