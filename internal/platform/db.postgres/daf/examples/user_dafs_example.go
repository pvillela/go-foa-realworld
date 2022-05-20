/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

var users = []model.User{
	{
		Username:     "pvillela",
		Email:        "foo@bar.com",
		PasswordHash: "dakfljads0fj",
		Bio:          util.PointerFromValue("I am me."),
		ImageLink:    nil,
	},
	{
		Username:     "joebloe",
		Email:        "joe@bloe.com",
		PasswordHash: "9zdakfljads0",
		Bio:          util.PointerFromValue("Famous person."),
		ImageLink:    util.PointerFromValue("https://myimage.com"),
	},
	{
		Username:     "johndoe",
		Email:        "johndoe@foo.com",
		PasswordHash: "09fs8asfoasi",
		Bio:          util.PointerFromValue("Average guy."),
		ImageLink:    util.PointerFromValue("https://johndooeimage.com"),
	},
}

func userDafsExample(ctx context.Context, ctxDb dbpgx.CtxPgx) {
	fmt.Println("********** userDafsExample **********\n")

	ctx, err := ctxDb.BeginTx(ctx)
	util.PanicOnError(err)
	//fmt.Println("ctx", ctx)

	for i, _ := range users {
		recCtx, err := daf.UserCreateDafI(ctx, &users[i])
		util.PanicOnError(err)
		_, _ = spew.Printf("user from Create: %v", users[i])
		fmt.Println("recCtx from Create:", recCtx, "\n")
	}

	userFromDb, recCtx, err := daf.UserGetByNameDafI(ctx, "pvillela")
	// Commented-out lines below were used to test forced error due to multiple rows returned
	//fmt.Println("Error classification:", dbpgx.ClassifyError(err))
	//uerr := errors.Unwrap(err)
	//fmt.Printf("Unwrapped error: %+v", uerr)
	util.PanicOnError(err)
	fmt.Println("UserGetByNameDaf:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx, "\n")

	userFromDb, recCtx, err = daf.UserGetByNameDafI(ctx, "xxx")
	fmt.Println("UserGetByNameDaf with invalid username")
	fmt.Println("Error:", err)
	fmt.Println("Error classification:", dbpgx.ClassifyError(err), "\n")
	uerr := errors.Unwrap(err)
	fmt.Printf("Unwrapped error: %+v\n\n", uerr)

	userFromDb, recCtx, err = daf.UserGetByEmailDafI(ctx, "foo@bar.com")
	util.PanicOnError(err)
	fmt.Println("UserGetByEmailDaf:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx, "\n")

	tx, err := dbpgx.GetCtxTx(ctx)
	util.PanicOnError(err)
	readManySql := "SELECT * FROM users"
	pwUsers, err := dbpgx.ReadMany[daf.PwUser](ctx, tx, readManySql, -1, -1)
	fmt.Println("pwUsers:", pwUsers, "\n")
	pwUserPtrs, err := dbpgx.ReadMany[*daf.PwUser](ctx, tx, readManySql, -1, -1)
	fmt.Println("pwUserPtrs:", pwUserPtrs, "\n")
	fmt.Println("*pwUserPtrs[0]:", *pwUserPtrs[0], "\n")

	ctx, err = ctxDb.Commit(ctx)
	util.PanicOnError(err)

	ctx, err = ctxDb.BeginTx(ctx)
	util.PanicOnError(err)

	user := users[0]
	user.ImageLink = util.PointerFromValue("https://xyz.com")
	recCtx, err = daf.UserUpdateDafI(ctx, user, recCtx)
	util.PanicOnError(err)
	fmt.Println("\nUserUpdateDaf:", user)
	fmt.Println("recCtx from Update:", recCtx, "\n")

	_, err = ctxDb.Commit(ctx)
	util.PanicOnError(err)
}
