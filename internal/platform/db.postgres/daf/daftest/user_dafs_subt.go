/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func UserDafsSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	errx.PanicOnError(err)

	ctx, err = ctxDb.BeginTx(ctx)
	errx.PanicOnError(err)
	defer ctxDb.DeferredRollback(ctx)

	{
		user := users[0]
		returned, recCtx, err := daf.UserGetByNameDafI(ctx, user.Username)
		errx.PanicOnError(err)
		util.Ignore(recCtx)
		//fmt.Println("UserGetByNameDaf:", userFromDb)
		//fmt.Println("recCtx from Read:", recCtx)

		expected := user
		assert.Equal(t, expected, returned)
	}

	{
		userFromDb, recCtx, err := daf.UserGetByNameDafI(ctx, "xxxxx")
		util.Ignore(userFromDb, recCtx)
		//fmt.Println("UserGetByNameDaf with invalid username")
		//fmt.Println("Error:", err)

		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		assert.Equal(t, expectedErrxKind, returnedErrxKind)
	}

	{
		user := users[1]
		returned, recCtx, err := daf.UserGetByEmailDafI(ctx, user.Email)
		errx.PanicOnError(err)
		util.Ignore(recCtx)
		//fmt.Println("UserGetByEmailDaf:", userFromDb)
		//fmt.Println("recCtx from Read:", recCtx)

		expected := user
		assert.Equal(t, expected, returned)
	}

	{
		tx, err := dbpgx.GetCtxTx(ctx)
		errx.PanicOnError(err)

		readManySql := "SELECT * FROM users"
		pwUsers, err := dbpgx.ReadMany[daf.PwUser](ctx, tx, readManySql, -1, -1)
		//fmt.Println("pwUsers:", pwUsers)

		returned := util.SliceMap(pwUsers, func(pw daf.PwUser) model.User {
			return pw.Entity
		})

		expected := users

		assert.ElementsMatch(t, expected, returned)
	}

	{
		userIdx := 0
		user := users[userIdx]
		recCtx := &recCtxUsers[userIdx]
		user.ImageLink = util.PointerFromValue("https://xyz.com")
		*recCtx, err = daf.UserUpdateDafI(ctx, user, *recCtx)
		errx.PanicOnError(err)
		//fmt.Println("\nUserUpdateDaf:", user)
		//fmt.Println("recCtx from Update:", recCtx)

		{
			var returned model.User
			returned, *recCtx, err = daf.UserGetByNameDafI(ctx, user.Username)
			errx.PanicOnError(err)
			//fmt.Println("UserGetByNameDaf:", userFromDb)
			//fmt.Println("recCtx from Read:", recCtx)

			expected := user
			assert.Equal(t, expected, returned)
		}
	}

	_, err = ctxDb.Commit(ctx)
	errx.PanicOnError(err)
	//fmt.Println("\nFinal commit for userDafsSubt")
}
