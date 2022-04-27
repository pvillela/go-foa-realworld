/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package newdaf

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// UserGetByNameDaf implements a stereotype instance of type
// fs.UserGetByNameDafT.
var UserGetByNameDaf fs.UserGetByNameDafT = func(
	ctx context.Context,
	userName string,
) (model.User, fs.RecCtxUser, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return model.User{}, fs.RecCtxUser{}, errx.ErrxOf(err)
	}
	rows, err := tx.Query(ctx, "SELECT * FROM users WHERE username = $1", userName)
	if err != nil {
		return model.User{}, fs.RecCtxUser{}, errx.ErrxOf(err)
	}
	pwUser := fs.PwUser{}
	err = pgxscan.ScanOne(&pwUser, rows)
	util.PanicOnError(err)
	return pwUser.Entity, pwUser.RecCtx, nil
}

// UserGetByEmailDaf implements a stereotype instance of type
// fs.UserGetByEmailDafT.
var UserGetByEmailDaf fs.UserGetByEmailDafT = func(
	ctx context.Context,
	email string,
) (model.User, fs.RecCtxUser, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return model.User{}, fs.RecCtxUser{}, errx.ErrxOf(err)
	}
	rows, err := tx.Query(ctx, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return model.User{}, fs.RecCtxUser{}, errx.ErrxOf(err)
	}
	pwUser := fs.PwUser{}
	err = pgxscan.ScanOne(&pwUser, rows)
	util.PanicOnError(err)
	return pwUser.Entity, pwUser.RecCtx, nil
}

// UserCreateDaf implements a stereotype instance of type
// fs.UserCreateDafT.
var UserCreateDaf fs.UserCreateDafT = func(
	ctx context.Context,
	user model.User,
) (fs.RecCtxUser, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return fs.RecCtxUser{}, errx.ErrxOf(err)
	}
	sql := `
	INSERT INTO users (username, email, password_hash, bio, image)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`
	args := []any{user.Username, user.Email, user.PasswordHash, user.Bio, user.ImageLink}
	row := tx.QueryRow(ctx, sql, args...)
	var recCtx fs.RecCtxUser
	err = row.Scan(&recCtx.Id, &recCtx.CreatedAt, &recCtx.UpdatedAt)
	if err != nil {
		return recCtx, errx.ErrxOf(err)
	}
	return recCtx, nil
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
