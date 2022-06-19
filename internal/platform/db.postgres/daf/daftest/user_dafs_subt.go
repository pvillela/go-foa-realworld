/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/cdb"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/sirupsen/logrus"
	"testing"

	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
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

	for i, _ := range users {
		user := users[i]
		recCtx, err := daf.UserCreateExplicitTxDafI(ctx, tx, &user)
		errx.PanicOnError(err)
		//_, _ = spew.Printf("user from Create: %v\n", user)
		logrus.Debug("user from Create:", user)
		logrus.Debug("recCtx from Create:", recCtx)

		mdb.UserUpsert2(user, recCtx)
	}
}

func userDafsSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	assert.NoError(t, err)

	_, err = cdb.WithTransaction(ctxDb, ctx, func(ctx context.Context) (types.Unit, error) {

		{
			tx, err := dbpgx.GetCtxTx(ctx)
			assert.NoError(t, err)
			setupUsers(ctx, tx)
		}

		{
			msg := "UserGetByNameDafI with valid username"

			username := username1

			retUser, retRecCtx, err := daf.UserGetByNameDafI(ctx, username)
			assert.NoError(t, err)
			util.Ignore(retRecCtx)
			//fmt.Println("UserGetByNameDaf:", userFromDb)
			//fmt.Println("recCtx from Read:", recCtx)

			expUser, expRecCtx := mdb.UserGet2(username)

			assert.Equal(t, expUser, retUser, msg+" - user")
			assert.Equal(t, expRecCtx, retRecCtx, msg+" = recCtx")
		}

		{
			msg := "UserGetByNameDafI with invalid username"

			username := "xxxxxx"

			_, _, err := daf.UserGetByNameDafI(ctx, username)
			returnedErrxKind := dbpgx.ClassifyError(err)
			expectedErrxKind := dbpgx.DbErrRecordNotFound

			assert.Equal(t, expectedErrxKind, returnedErrxKind, msg)
		}

		{
			msg := "UserGetByEmailDafI with valid email"

			username := username2

			expUser, expRecCtx := mdb.UserGet2(username)

			retUser, retRecCtx, err := daf.UserGetByEmailDafI(ctx, expUser.Email)
			assert.NoError(t, err)
			//fmt.Println("UserGetByEmailDaf:", userFromDb)
			//fmt.Println("recCtx from Read:", recCtx)

			assert.Equal(t, expUser, retUser, msg+" - user")
			assert.Equal(t, expRecCtx, retRecCtx, msg+" - recCtx")
		}

		{
			msg := "UserGetByEmailDafI with invalid email"

			email := "xxxxxx@xxx.xx"

			_, _, err := daf.UserGetByEmailDafI(ctx, email)
			//fmt.Println("UserGetByNameDaf with invalid username")
			//fmt.Println("Error:", err)

			returnedErrxKind := dbpgx.ClassifyError(err)
			expectedErrxKind := dbpgx.DbErrRecordNotFound

			assert.Equal(t, expectedErrxKind, returnedErrxKind, msg)
		}

		{
			msg := "Retrieve all users"

			tx, err := dbpgx.GetCtxTx(ctx)
			assert.NoError(t, err)

			readManySql := "SELECT * FROM users"
			pwUsers, err := dbpgx.ReadMany[daf.PwUser](ctx, tx, readManySql, -1, -1)
			//fmt.Println("pwUsers:", pwUsers)

			returned := util.SliceMap(pwUsers, func(pw daf.PwUser) model.User {
				return pw.Entity
			})

			expected, _ := mdb.UserGet2All()

			assert.ElementsMatch(t, expected, returned, msg)
		}

		{
			msg := "UserUpdateDafI - image"

			username := username1

			user, recCtx := mdb.UserGet2(username)
			user.ImageLink = "https://xyz.com"

			updRecCtx, err := daf.UserUpdateDafI(ctx, user, recCtx)
			assert.NoError(t, err)
			//fmt.Println("\nUserUpdateDaf:", user)
			//fmt.Println("recCtx from Update:", recCtx)

			assert.Equal(t, recCtx.CreatedAt, updRecCtx.CreatedAt, msg+" - recCtx.CreatedAt must be equal to updRecCtx.CreatedAt")
			assert.NotEqual(t, recCtx.UpdatedAt, updRecCtx.UpdatedAt, msg+" - recCtx.UpdatedAt must be different from updRecCtx.UpdatedAt")

			// Sync in-memory database
			mdb.UserUpsert2(user, updRecCtx)

			{
				retUser, retRecCtx, err := daf.UserGetByNameDafI(ctx, user.Username)
				assert.NoError(t, err)
				//fmt.Println("UserGetByNameDaf:", userFromDb)
				//fmt.Println("recCtx from Read:", recCtx)

				assert.Equal(t, user, retUser, msg+" - user must be equal to retUser")
				assert.Equal(t, updRecCtx, retRecCtx, msg+" - reqCtx must be equal to retReqCtx")
			}
		}

		{
			msg := "UserUpdateDafI - bio"

			username := username2

			user, recCtx := mdb.UserGet2(username)
			user.Bio = util.PointerFromValue("I'm a really famous person.")

			updRecCtx, err := daf.UserUpdateDafI(ctx, user, recCtx)
			assert.NoError(t, err)
			//fmt.Println("\nUserUpdateDaf:", user)
			//fmt.Println("recCtx from Update:", recCtx)

			assert.Equal(t, recCtx.CreatedAt, updRecCtx.CreatedAt, msg+" - recCtx.CreatedAt must be equal to updRecCtx.CreatedAt")
			assert.NotEqual(t, recCtx.UpdatedAt, updRecCtx.UpdatedAt, msg+" - recCtx.UpdatedAt must be different from updRecCtx.UpdatedAt")

			// Sync in-memory database
			mdb.UserUpsert2(user, updRecCtx)

			{
				retUser, retRecCtx, err := daf.UserGetByNameDafI(ctx, user.Username)
				assert.NoError(t, err)
				//fmt.Println("UserGetByNameDaf:", userFromDb)
				//fmt.Println("recCtx from Read:", recCtx)

				assert.Equal(t, user, retUser, msg+" - user must be equal to retUser")
				assert.Equal(t, updRecCtx, retRecCtx, msg+" - reqCtx must be equal to retReqCtx")
			}
		}

		return types.UnitV, err
	})

	assert.NoError(t, err)
}
