/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

func userDafsExample(ctx context.Context, ctxDb dbpgx.CtxPgx) {
	users := []model.User{
		{
			Username:     "pvillela",
			Email:        "foo@bar.com",
			PasswordHash: "dakfljads0fj",
			Bio:          util.PointerFromValue("I am me."),
			ImageLink:    "",
		},
		{
			Username:     "joebloe",
			Email:        "joe@bloe.com",
			PasswordHash: "9zdakfljads0",
			Bio:          util.PointerFromValue("Famous person."),
			ImageLink:    "https://myimage.com",
		},
	}

	ctx, err := ctxDb.BeginTx(ctx)
	util.PanicOnError(err)
	//fmt.Println("ctx", ctx)

	for i, _ := range users {
		recCtx, err := daf.UserCreateDaf(ctx, &users[i])
		util.PanicOnError(err)
		fmt.Println("recCtx from Create:", recCtx)
	}

	userFromDb, recCtx, err := daf.UserGetByNameDaf(ctx, "pvillela")
	util.PanicOnError(err)
	fmt.Println("\nUserGetByNameDaf:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx)

	userFromDb, recCtx, err = daf.UserGetByEmailDaf(ctx, "foo@bar.com")
	util.PanicOnError(err)
	fmt.Println("\nUserGetByEmailDaf:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx)

	tx, err := dbpgx.GetCtxTx(ctx)
	util.PanicOnError(err)
	readManySql := "SELECT * FROM users"
	pwUsers, err := dbpgx.ReadMany[daf.PwUser](ctx, tx, readManySql, -1, -1)
	fmt.Println("\npwUsers:", pwUsers)
	pwUserPtrs, err := dbpgx.ReadMany[*daf.PwUser](ctx, tx, readManySql, -1, -1)
	fmt.Println("pwUserPtrs:", pwUserPtrs)
	fmt.Println("*pwUserPtrs[0]:", *pwUserPtrs[0])

	ctx, err = ctxDb.Commit(ctx)
	util.PanicOnError(err)

	ctx, err = ctxDb.BeginTx(ctx)
	util.PanicOnError(err)

	user := users[0]
	user.ImageLink = "https://xyz.com"
	recCtx, err = daf.UserUpdateDaf(ctx, user, recCtx)
	util.PanicOnError(err)
	fmt.Println("\nUserUpdateDaf:", user)
	fmt.Println("recCtx from Update:", recCtx)

	_, err = ctxDb.Commit(ctx)
	util.PanicOnError(err)
}
