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
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/sirupsen/logrus"
	"testing"

	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/stretchr/testify/assert"
)

func setupUsers(ctx context.Context, tx pgx.Tx) {
	var users = []model.User{
		{
			Username:     "pvillela",
			Email:        "foo@bar.com",
			PasswordHash: "dakfljads0fj",
			PasswordSalt: "2af8d0b50a",
			Bio:          util.PointerFromValue("I am me."),
			ImageLink:    nil,
		},
		{
			Username:     "joebloe",
			Email:        "joe@bloe.com",
			PasswordHash: "9zdakfljads0",
			PasswordSalt: "3ba9e9c611",
			Bio:          util.PointerFromValue("Famous person."),
			ImageLink:    util.PointerFromValue("https://myimage.com"),
		},
		{
			Username:     "johndoe",
			Email:        "johndoe@foo.com",
			PasswordHash: "09fs8asfoasi",
			PasswordSalt: "0000000000",
			Bio:          util.PointerFromValue("Average guy."),
			ImageLink:    util.PointerFromValue("https://johndooeimage.com"),
		},
	}

	for i, _ := range users {
		user := users[i]
		recCtx, err := daf.UserCreateExplicitTxDafI(ctx, tx, &user)
		errx.PanicOnError(err)
		//_, _ = spew.Printf("user from Create: %v\n", user)
		logrus.Debug("user from Create:", user)
		logrus.Debug("recCtx from Create:", recCtx)

		mdb.UserUpsert(user.Username, user, recCtx)
	}
}

func userDafsSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	errx.PanicOnError(err)

	_, err = cdb.WithTransaction(ctxDb, ctx, func(ctx context.Context) (types.Unit, error) {

		{
			tx, err := dbpgx.GetCtxTx(ctx)
			errx.PanicOnError(err)
			setupUsers(ctx, tx)
		}

		{
			msg := "UserGetByNameDafI with valid username"

			username := "pvillela"

			retUser, retRecCtx, err := daf.UserGetByNameDafI(ctx, username)
			errx.PanicOnError(err)
			util.Ignore(retRecCtx)
			//fmt.Println("UserGetByNameDaf:", userFromDb)
			//fmt.Println("recCtx from Read:", recCtx)

			expUser, expRecCtx := mdb.UserGet(username)

			assert.Equal(t, expUser, retUser, msg+" - user")
			assert.Equal(t, expRecCtx, retRecCtx, msg+" = recCtx")
		}

		{
			msg := "UserGetByNameDafI with invalid username"

			username := "xxxxxx"

			_, _, err := daf.UserGetByNameDafI(ctx, username)
			//fmt.Println("UserGetByNameDaf with invalid username")
			//fmt.Println("Error:", err)

			returnedErrxKind := dbpgx.ClassifyError(err)
			expectedErrxKind := dbpgx.DbErrRecordNotFound

			assert.Equal(t, expectedErrxKind, returnedErrxKind, msg)
		}

		{
			msg := "UserGetByEmailDafI with valid email"

			username := "joebloe"

			expUser, expRecCtx := mdb.UserGet(username)

			retUser, retRecCtx, err := daf.UserGetByEmailDafI(ctx, expUser.Email)
			errx.PanicOnError(err)
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
			errx.PanicOnError(err)

			readManySql := "SELECT * FROM users"
			pwUsers, err := dbpgx.ReadMany[daf.PwUser](ctx, tx, readManySql, -1, -1)
			//fmt.Println("pwUsers:", pwUsers)

			returned := util.SliceMap(pwUsers, func(pw daf.PwUser) model.User {
				return pw.Entity
			})

			expected, _ := mdb.Users()

			assert.ElementsMatch(t, expected, returned, msg)
		}

		{
			msg := "UserUpdateDafI - image"

			username := "pvillela"

			user, recCtx := mdb.UserGet(username)
			user.ImageLink = util.PointerFromValue("https://xyz.com")

			updRecCtx, err := daf.UserUpdateDafI(ctx, user, recCtx)
			errx.PanicOnError(err)
			//fmt.Println("\nUserUpdateDaf:", user)
			//fmt.Println("recCtx from Update:", recCtx)

			assert.Equal(t, recCtx.CreatedAt, updRecCtx.CreatedAt, msg+" - recCtx.CreatedAt must be equal to updRecCtx.CreatedAt")
			assert.NotEqual(t, recCtx.UpdatedAt, updRecCtx.UpdatedAt, msg+" - recCtx.UpdatedAt must be different from updRecCtx.UpdatedAt")

			// Sync in-memory database
			mdb.UserUpsert(username, user, updRecCtx)

			{
				retUser, retRecCtx, err := daf.UserGetByNameDafI(ctx, user.Username)
				errx.PanicOnError(err)
				//fmt.Println("UserGetByNameDaf:", userFromDb)
				//fmt.Println("recCtx from Read:", recCtx)

				assert.Equal(t, user, retUser, msg+" - user must be equal to retUser")
				assert.Equal(t, updRecCtx, retRecCtx, msg+" - reqCtx must be equal to retReqCtx")
			}
		}

		{
			msg := "UserUpdateDafI - bio"

			username := "joebloe"

			user, recCtx := mdb.UserGet(username)
			user.Bio = util.PointerFromValue("I'm a really famous person.")

			updRecCtx, err := daf.UserUpdateDafI(ctx, user, recCtx)
			errx.PanicOnError(err)
			//fmt.Println("\nUserUpdateDaf:", user)
			//fmt.Println("recCtx from Update:", recCtx)

			assert.Equal(t, recCtx.CreatedAt, updRecCtx.CreatedAt, msg+" - recCtx.CreatedAt must be equal to updRecCtx.CreatedAt")
			assert.NotEqual(t, recCtx.UpdatedAt, updRecCtx.UpdatedAt, msg+" - recCtx.UpdatedAt must be different from updRecCtx.UpdatedAt")

			// Sync in-memory database
			mdb.UserUpsert(username, user, updRecCtx)

			{
				retUser, retRecCtx, err := daf.UserGetByNameDafI(ctx, user.Username)
				errx.PanicOnError(err)
				//fmt.Println("UserGetByNameDaf:", userFromDb)
				//fmt.Println("recCtx from Read:", recCtx)

				assert.Equal(t, user, retUser, msg+" - user must be equal to retUser")
				assert.Equal(t, updRecCtx, retRecCtx, msg+" - reqCtx must be equal to retReqCtx")
			}
		}

		return types.UnitV, err
	})

	errx.PanicOnError(err)
}
