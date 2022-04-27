/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/newdaf"
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

	ctx, err = ctxDb.Begin(ctx)
	util.PanicOnError(err)
	//fmt.Println("ctx", ctx)

	recCtx, err := newdaf.UserCreateDaf(ctx, user)
	util.PanicOnError(err)
	fmt.Println("recCtx from Create:", recCtx)

	userFromDb, recCtx, err := newdaf.UserGetByNameDaf(ctx, "pvillela")
	util.PanicOnError(err)
	fmt.Println("\nUserGetByNameDaf:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx)

	userFromDb, recCtx, err = newdaf.UserGetByEmailDaf(ctx, "foo@bar.com")
	util.PanicOnError(err)
	fmt.Println("\nUserGetByEmailDaf:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx)

	//err = ctxDb.Commit(ctx)
	//util.PanicOnError(err)
}
