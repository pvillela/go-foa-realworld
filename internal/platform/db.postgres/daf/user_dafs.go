/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"

	"github.com/georgysavva/scany/pgxscan"
)

// UserGetByNameDaf implements a stereotype instance of type
// fs.UserGetByNameDafT.
var UserGetByNameDaf fs.UserGetByNameDafT = func(
	ctx context.Context,
	userName string,
) (model.User, fs.RecCtxUser, error) {
	conn, err := dbpgx.GetCtxConn(ctx)
	if err != nil {
		return model.User{}, fs.RecCtxUser{}, err
	}
	rows, err := conn.Query(ctx, "SELECT * FROM users WHERE user_name = $1", userName)
	if err != nil {
		return model.User{}, fs.RecCtxUser{}, err
	}
	pwUser := fs.PwUser{}
	err = pgxscan.ScanOne(&pwUser, rows)
	util.PanicOnError(err)
	return pwUser.Entity, pwUser.RecCtx, nil
}

//// UserGetByEmailDaf implements a stereotype instance of type
//// fs.UserGetByEmailDafT.
//var UserGetByEmailDaf fs.UserGetByEmailDafT = func(
//	ctx context.Context,
//	email string,
//) (model.User, fs.RecCtxUser, error) {
//	conn := GetCtxConn(ctx)
//	rows, err := conn.Query(ctx, "SELECT * FROM users WHERE email = $1", email)
//	util.PanicOnError(err)
//	user := model.User{}
//	err = pgxscan.ScanOne(&user, rows)
//	util.PanicOnError(err)
//	return user, fs.RecCtxUser{}, nil // TODO: RecCtx is empty; decide if it needs a non-empty value or should be discarded
//}

// UserCreateDaf implements a stereotype instance of type
// fs.UserCreateDafT.
var UserCreateDaf fs.UserCreateDafT = func(
	ctx context.Context,
	user model.User,
) (fs.RecCtxUser, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return fs.RecCtxUser{}, err
	}
	_, err = tx.Exec(ctx, "INSERT INTO users VALUES ($1, $2, $3, $4, $5)",
		user.Username, user.Email, user.PasswordHash, user.Bio,
		user.ImageLink)
	if err != nil {
		return fs.RecCtxUser{}, err
	}
	return fs.RecCtxUser{}, nil // TODO: return proper RecCtxUser
}

//// UserUpdateDafC is the function that constructs a stereotype instance of type
//// fs.UserUpdateDafT.
//func UserUpdateDafC(
//	userDb mapdb.MapDb,
//) fs.UserUpdateDafT {
//	return func(user model.User, recCtx fs.RecCtxUser, txn db.Txn) (fs.RecCtxUser, error) {
//		if userByEmail, _, err := getByEmail(userDb, user.Email); err == nil && userByEmail.Name != user.Username {
//			return fs.RecCtxUser{}, fs.ErrDuplicateUserEmail.Make(nil, user.Email)
//		}
//
//		pw := fs.PwUser{RecCtx: recCtx, Entity: user}
//		err := userDb.Update(user.Username, pw, txn)
//		if errx.KindOf(err) == mapdb.ErrRecordNotFound {
//			return fs.RecCtxUser{}, fs.ErrUserNameNotFound.Make(err, user.Username)
//		}
//		if err != nil {
//			return fs.RecCtxUser{}, err // this can only be a transaction error
//		}
//
//		return recCtx, nil
//	}
//}