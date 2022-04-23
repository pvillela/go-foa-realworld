/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"

	"github.com/georgysavva/scany/pgxscan"
)

// UserGetByNameDaf implements a stereotype instance of type
// fs.UserGetByNameDafT.
var UserGetByNameDaf fs.UserGetByNameDafT = func(
	reqCtx context.Context,
	userName string,
) (model.User, fs.RecCtxUser, error) {
	conn := GetCtxConn(reqCtx)
	rows, err := conn.Query(reqCtx, "SELECT * FROM users WHERE user_name = $1", userName)
	util.PanicOnError(err)
	user := model.User{}
	err = pgxscan.ScanOne(&user, rows)
	util.PanicOnError(err)
	return user, fs.RecCtxUser{}, nil // TODO: RecCtx is empty; decide if it needs a non-empty value or should be discarded
}

// UserGetByEmailDaf implements a stereotype instance of type
// fs.UserGetByEmailDafT.
var UserGetByEmailDaf fs.UserGetByEmailDafT = func(
	reqCtx context.Context,
	email string,
) (model.User, fs.RecCtxUser, error) {
	conn := GetCtxConn(reqCtx)
	rows, err := conn.Query(reqCtx, "SELECT * FROM users WHERE email = $1", email)
	util.PanicOnError(err)
	user := model.User{}
	err = pgxscan.ScanOne(&user, rows)
	util.PanicOnError(err)
	return user, fs.RecCtxUser{}, nil // TODO: RecCtx is empty; decide if it needs a non-empty value or should be discarded
}

// UserCreateDaf implements a stereotype instance of type
// fs.UserCreateDafT.
var UserCreateDaf fs.UserCreateDafT = func(
	reqCtx context.Context,
	user model.User,
	txn db.Txn,
) (fs.RecCtxUser, error) {
	conn := GetCtxConn(reqCtx)
	tx, err := conn.Begin(reqCtx)
	util.PanicOnError(err)
	defer tx.Rollback(reqCtx)
	_, err = tx.Exec(reqCtx, "INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		user.Username, user.Email, user.IsTempPassword, user.PasswordHash, user.PasswordSalt, user.Bio,
		user.ImageLink, user.Followees, user.NumFollowers, user.CreatedAt, user.UpdatedAt)
	util.PanicOnError(err)
	_, err = conn.Exec(reqCtx, "UserCreate", "foo")
}

// UserUpdateDafC is the function that constructs a stereotype instance of type
// fs.UserUpdateDafT.
func UserUpdateDafC(
	userDb mapdb.MapDb,
) fs.UserUpdateDafT {
	return func(user model.User, recCtx fs.RecCtxUser, txn db.Txn) (fs.RecCtxUser, error) {
		if userByEmail, _, err := getByEmail(userDb, user.Email); err == nil && userByEmail.Name != user.Username {
			return fs.RecCtxUser{}, fs.ErrDuplicateUserEmail.Make(nil, user.Email)
		}

		pw := fs.PwUser{RecCtx: recCtx, Entity: user}
		err := userDb.Update(user.Username, pw, txn)
		if errx.KindOf(err) == mapdb.ErrRecordNotFound {
			return fs.RecCtxUser{}, fs.ErrUserNameNotFound.Make(err, user.Username)
		}
		if err != nil {
			return fs.RecCtxUser{}, err // this can only be a transaction error
		}

		return recCtx, nil
	}
}
