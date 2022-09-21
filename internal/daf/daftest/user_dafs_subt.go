/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx/dbpgxtest"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/sirupsen/logrus"
	"testing"

	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/stretchr/testify/assert"
)

const (
	username1 = "pvillela"
	username2 = "joebloe"
	username3 = "johndoe"
)

func setupUsers(ctx context.Context, tx pgx.Tx) {
	var users = []model.User{
		{
			Username:     username1,
			Email:        "foo@bar.com",
			PasswordHash: "dakfljads0fj",
			PasswordSalt: "2af8d0b50a",
			Bio:          util.PointerOf("I am me."),
			ImageLink:    "",
		},
		{
			Username:     username2,
			Email:        "joe@bloe.com",
			PasswordHash: "9zdakfljads0",
			PasswordSalt: "3ba9e9c611",
			Bio:          util.PointerOf("Famous person."),
			ImageLink:    "https://myimage.com",
		},
		{
			Username:     username3,
			Email:        "johndoe@foo.com",
			PasswordHash: "09fs8asfoasi",
			PasswordSalt: "0000000000",
			Bio:          util.PointerOf("Average guy."),
			ImageLink:    "https://johndooeimage.com",
		},
	}

	for i, _ := range users {
		user := users[i]
		err := daf.UserCreateDaf(ctx, tx, &user)
		errx.PanicOnError(err)
		//_, _ = spew.Printf("user from Create: %v\n", user)
		logrus.Debug("user from Create:", user)
		logrus.Debug("user from Create:")

		mdb.UserUpsert(user)
	}
}

var userDafsSubt = dbpgxtest.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {

	{
		setupUsers(ctx, tx)
	}

	{
		msg := "UserGetByNameDaf with valid username"

		username := username1

		retUser, err := daf.UserGetByNameDaf(ctx, tx, username)
		assert.NoError(t, err)
		//fmt.Println("UserGetByNameDaf:", userFromDb)
		//fmt.Println("user from Read:")

		expUser := mdb.UserGetByName(username)

		assert.Equal(t, expUser, retUser, msg+" - user")
	}

	{
		msg := "UserGetByNameDaf with invalid username"

		username := "xxxxxx"

		_, err := daf.UserGetByNameDaf(ctx, tx, username)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg)
	}

	{
		msg := "UserGetByEmailDaf with valid email"

		username := username2

		expUser := mdb.UserGetByName(username)

		retUser, err := daf.UserGetByEmailDaf(ctx, tx, expUser.Email)
		assert.NoError(t, err)
		//fmt.Println("UserGetByEmailDaf:", userFromDb)
		//fmt.Println("user from Read:")

		assert.Equal(t, expUser, retUser, msg+" - user")
	}

	{
		msg := "UserGetByEmailDaf with invalid email"

		email := "xxxxxx@xxx.xx"

		_, err := daf.UserGetByEmailDaf(ctx, tx, email)
		//fmt.Println("UserGetByNameDaf with invalid username")
		//fmt.Println("Error:", err)

		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg)
	}

	{
		msg := "Retrieve all users"

		readManySql := "SELECT * FROM users"
		returned, err := dbpgx.ReadMany[model.User](ctx, tx, readManySql, -1, -1)
		assert.NoError(t, err)
		//fmt.Println("pwUsers:", pwUsers)

		expected := mdb.UserGetAll()

		assert.ElementsMatch(t, expected, returned, msg)
	}

	{
		msg := "UserUpdateDaf - image"

		username := username1

		user := mdb.UserGetByName(username)
		user.ImageLink = "https://xyz.com"

		updUser := user
		err := daf.UserUpdateDaf(ctx, tx, &updUser)
		assert.NoError(t, err)
		//fmt.Println("\nUserUpdateDaf:", user)
		//fmt.Println("updUser from Update:", updUser)

		assert.Equal(t, user.CreatedAt, updUser.CreatedAt, msg+" - user.CreatedAt must be equal to updUser.CreatedAt")
		assert.NotEqual(t, user.UpdatedAt, updUser.UpdatedAt, msg+" - user.UpdatedAt must be different from updUser.UpdatedAt")

		// Sync in-memory database
		mdb.UserUpsert(updUser)

		{
			retUser, err := daf.UserGetByNameDaf(ctx, tx, user.Username)
			assert.NoError(t, err)
			//fmt.Println("UserGetByNameDaf:", userFromDb)
			//fmt.Println("user from Read:", retUser)

			assert.Equal(t, updUser, retUser, msg+" - updUser must be equal to retUser")
		}
	}

	{
		msg := "UserUpdateDaf - bio"

		username := username2

		user := mdb.UserGetByName(username)
		user.Bio = util.PointerOf("I'm a really famous person.")

		updUser := user
		err := daf.UserUpdateDaf(ctx, tx, &updUser)
		assert.NoError(t, err)
		//fmt.Println("\nUserUpdateDaf:", user)
		//fmt.Println("user from Update:", updUser)

		assert.Equal(t, user.CreatedAt, updUser.CreatedAt, msg+" - user.CreatedAt must be equal to updUser.CreatedAt")
		assert.NotEqual(t, user.UpdatedAt, updUser.UpdatedAt, msg+" - user.UpdatedAt must be different from updUser.UpdatedAt")

		// Sync in-memory database
		mdb.UserUpsert(updUser)

		{
			retUser, err := daf.UserGetByNameDaf(ctx, tx, user.Username)
			assert.NoError(t, err)
			//fmt.Println("UserGetByNameDaf:", userFromDb)
			//fmt.Println("user from Read:", retUser)

			assert.Equal(t, updUser, retUser, msg+" - updUser must be equal to retUser")
		}
	}
})
