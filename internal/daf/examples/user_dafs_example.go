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
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

const (
	username1 = "pvillela"
	username2 = "joebloe"
	username3 = "johndoe"
)

var users = []model.User{
	{
		Username:     username1,
		Email:        "foo@bar.com",
		PasswordHash: "dakfljads0fj",
		PasswordSalt: "2af8d0b50a",
		Bio:          util.PointerFromValue("I am me."),
		ImageLink:    "",
	},
	{
		Username:     username2,
		Email:        "joe@bloe.com",
		PasswordHash: "9zdakfljads0",
		PasswordSalt: "3ba9e9c611",
		Bio:          util.PointerFromValue("Famous person."),
		ImageLink:    "https://myimage.com",
	},
	{
		Username:     username3,
		Email:        "johndoe@foo.com",
		PasswordHash: "09fs8asfoasi",
		PasswordSalt: "0000000000",
		Bio:          util.PointerFromValue("Average guy."),
		ImageLink:    "https://johndooeimage.com",
	},
}

func userDafsExample(ctx context.Context, db dbpgx.Db) {
	fmt.Println("********** userDafsExample **********\n")

	tx, err := db.BeginTx(ctx)
	errx.PanicOnError(err)
	//fmt.Println("ctx", ctx)

	for i, _ := range users {
		err := daf.UserCreateDaf(ctx, tx, &users[i])
		errx.PanicOnError(err)
		_, _ = spew.Printf("user from Create: %v\n", users[i])
	}

	userFromDb, err := daf.UserGetByNameDaf(ctx, tx, username1)
	// Commented-out lines below were used to daftest forced error due to multiple rows returned
	//fmt.Println("Error classification:", dbpgx.ClassifyError(err))
	//uerr := errors.Unwrap(err)
	//fmt.Printf("Unwrapped error: %+v", uerr)
	errx.PanicOnError(err)
	fmt.Println("UserGetByNameDaf:", userFromDb)

	userFromDb, err = daf.UserGetByNameDaf(ctx, tx, "daftest")
	fmt.Println("UserGetByNameDaf with invalid username")
	fmt.Println("Error:", err)
	fmt.Println("Error classification:", dbpgx.ClassifyError(err), "\n")
	uerr := errors.Unwrap(err)
	fmt.Printf("Unwrapped error: %+v\n\n", uerr)

	userFromDb, err = daf.UserGetByEmailDaf(ctx, tx, "foo@bar.com")
	errx.PanicOnError(err)
	fmt.Println("UserGetByEmailDaf:", userFromDb)

	readManySql := "SELECT * FROM users"
	users, err := dbpgx.ReadMany[model.User](ctx, tx, readManySql, -1, -1)
	fmt.Println("users:", users, "\n")
	userPtrs, err := dbpgx.ReadMany[*model.User](ctx, tx, readManySql, -1, -1)
	fmt.Println("userPtrs:", userPtrs, "\n")
	fmt.Println("*userPtrs[0]:", *userPtrs[0], "\n")

	err = tx.Commit(ctx)
	errx.PanicOnError(err)

	tx, err = db.BeginTx(ctx)
	errx.PanicOnError(err)

	user := users[0]
	user.ImageLink = "https://xyz.com"
	err = daf.UserUpdateDaf(ctx, tx, &user)
	errx.PanicOnError(err)
	fmt.Println("\nUserUpdateDaf:", user)

	err = tx.Commit(ctx)
	errx.PanicOnError(err)
}
