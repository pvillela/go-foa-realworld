/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"errors"
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"testing"
)

func UserDafsSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	errx.PanicOnError(err)

	ctx, err = ctxDb.BeginTx(ctx)
	errx.PanicOnError(err)
	//fmt.Println("ctx", ctx)

	userFromDb, recCtx, err := daf.UserGetByNameDafI(ctx, "pvillela")
	// Commented-out lines below were used to daftest forced error due to multiple rows returned
	//fmt.Println("Error classification:", dbpgx.ClassifyError(err))
	//uerr := errors.Unwrap(err)
	//fmt.Printf("Unwrapped error: %+v", uerr)
	errx.PanicOnError(err)
	fmt.Println("UserGetByNameDaf:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx)

	userFromDb, recCtx, err = daf.UserGetByNameDafI(ctx, "daftest")
	fmt.Println("UserGetByNameDaf with invalid username")
	fmt.Println("Error:", err)
	fmt.Println("Error classification:", dbpgx.ClassifyError(err))
	uerr := errors.Unwrap(err)
	fmt.Printf("Unwrapped error: %+v\n\n", uerr)

	userFromDb, recCtx, err = daf.UserGetByEmailDafI(ctx, "foo@bar.com")
	errx.PanicOnError(err)
	fmt.Println("UserGetByEmailDaf:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx)

	tx, err := dbpgx.GetCtxTx(ctx)
	errx.PanicOnError(err)
	readManySql := "SELECT * FROM users"
	pwUsers, err := dbpgx.ReadMany[daf.PwUser](ctx, tx, readManySql, -1, -1)
	fmt.Println("pwUsers:", pwUsers)
	pwUserPtrs, err := dbpgx.ReadMany[*daf.PwUser](ctx, tx, readManySql, -1, -1)
	fmt.Println("pwUserPtrs:", pwUserPtrs)
	fmt.Println("*pwUserPtrs[0]:", *pwUserPtrs[0])

	ctx, err = ctxDb.Commit(ctx)
	errx.PanicOnError(err)

	ctx, err = ctxDb.BeginTx(ctx)
	errx.PanicOnError(err)

	user := users[0]
	user.ImageLink = util.PointerFromValue("https://xyz.com")
	recCtx, err = daf.UserUpdateDafI(ctx, user, recCtx)
	errx.PanicOnError(err)
	fmt.Println("\nUserUpdateDaf:", user)
	fmt.Println("recCtx from Update:", recCtx)

	_, err = ctxDb.Commit(ctx)
	errx.PanicOnError(err)
	fmt.Println("\nFinal commit for userDafsSubt")
}

func UserDafsSubt1(db dbpgx.Db, ctx context.Context) {
	////defer errx.PanicLog(log.Fatal)
	////
	////log.SetLevel(log.DebugLevel)
	//////var arr []any
	//////fmt.Println(arr[0])
	////
	////ctx = context.Background()
	//
	//connStr := "postgres://testuser:testpassword@localhost:9999/testdb?sslmode=disable"
	//pool, err := pgxpool.Connect(ctx, connStr)
	//errx.PanicOnError(err)
	//fmt.Println("pool:", pool)
	//
	//ctxDb := dbpgx.CtxPgx{pool}
	//ctx, err = ctxDb.SetPool(ctx)
	//errx.PanicOnError(err)

	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	errx.PanicOnError(err)

	ctx, err = ctxDb.BeginTx(ctx)
	errx.PanicOnError(err)
	//fmt.Println("ctx", ctx)

	userFromDb, recCtx, err := daf.UserGetByNameDafI(ctx, "pvillela")
	// Commented-out lines below were used to daftest forced error due to multiple rows returned
	//fmt.Println("Error classification:", dbpgx.ClassifyError(err))
	//uerr := errors.Unwrap(err)
	//fmt.Printf("Unwrapped error: %+v", uerr)
	errx.PanicOnError(err)
	fmt.Println("UserGetByNameDaf:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx)

	//userFromDb, recCtx, err = daf.UserGetByNameDafI(ctx, "daftest")
	//fmt.Println("UserGetByNameDaf with invalid username")
	//fmt.Println("Error:", err)
	//fmt.Println("Error classification:", dbpgx.ClassifyError(err))
	//uerr := errors.Unwrap(err)
	//fmt.Printf("Unwrapped error: %+v\n\n", uerr)
	//
	//userFromDb, recCtx, err = daf.UserGetByEmailDafI(ctx, "foo@bar.com")
	//errx.PanicOnError(err)
	//fmt.Println("UserGetByEmailDaf:", userFromDb)
	//fmt.Println("recCtx from Read:", recCtx)
	//
	//tx, err := dbpgx.GetCtxTx(ctx)
	//errx.PanicOnError(err)
	//readManySql := "SELECT * FROM users"
	//pwUsers, err := dbpgx.ReadMany[daf.PwUser](ctx, tx, readManySql, -1, -1)
	//fmt.Println("pwUsers:", pwUsers)
	//pwUserPtrs, err := dbpgx.ReadMany[*daf.PwUser](ctx, tx, readManySql, -1, -1)
	//fmt.Println("pwUserPtrs:", pwUserPtrs)
	//fmt.Println("*pwUserPtrs[0]:", *pwUserPtrs[0])
	//
	//ctx, err = ctxDb.Commit(ctx)
	//errx.PanicOnError(err)
	//
	//ctx, err = ctxDb.BeginTx(ctx)
	//errx.PanicOnError(err)
	//
	//user := users[0]
	//user.ImageLink = util.PointerFromValue("https://xyz.com")
	//recCtx, err = daf.UserUpdateDafI(ctx, user, recCtx)
	//errx.PanicOnError(err)
	//fmt.Println("\nUserUpdateDaf:", user)
	//fmt.Println("recCtx from Update:", recCtx)

	_, err = ctxDb.Commit(ctx)
	errx.PanicOnError(err)
	fmt.Println("\nFinal commit for userDafsSubt")
}
