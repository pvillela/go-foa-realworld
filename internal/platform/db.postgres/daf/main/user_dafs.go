/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

func main() {
	defer util.PanicLog()
	myBio := "I am me."
	user := model.User{
		Username:     "pvillela",
		Email:        "foo@bar.com",
		PasswordHash: "dakfljads0fj",
		Bio:          &myBio,
		ImageLink:    "",
	}

	ctx := context.Background()

	connStr := "postgres://testuser:testpassword@localhost:9999/testdb?sslmode=disable"
	pool, err := pgxpool.Connect(ctx, connStr)
	util.PanicOnError(err)
	ctxConn := dbpgx.CtxPgx{pool}
	ctx, err = ctxConn.SetPool(ctx)
	util.PanicOnError(err)
	ctx, err = ctxConn.Begin(ctx)
	util.PanicOnError(err)
	fmt.Println("ctx", ctx)

	recCtx, err := UserCreateDaf(ctx, user)
	util.PanicOnError(err)
	fmt.Println("recCtx from Create:", recCtx)

	userFromDb, recCtx, err := UserGetByNameDaf(ctx, "pvillela")
	util.PanicOnError(err)
	fmt.Println("userFromDb:", userFromDb)
	fmt.Println("recCtx from Read:", recCtx)
}

// UserGetByNameDaf implements a stereotype instance of type
// fs.UserGetByNameDafT.
var UserGetByNameDaf fs.UserGetByNameDafT = func(
	ctx context.Context,
	userName string,
) (model.User, fs.RecCtxUser, error) {
	conn, err := dbpgx.GetCtxConn(ctx)
	if err != nil {
		return model.User{}, fs.RecCtxUser{}, errx.ErrxOf(err)
	}
	rows, err := conn.Query(ctx, "SELECT * FROM users WHERE username = $1", userName)
	if err != nil {
		return model.User{}, fs.RecCtxUser{}, errx.ErrxOf(err)
	}
	//pwUser := fs.PwUser{}
	user := model.User{}
	err = pgxscan.ScanOne(&user, rows)
	util.PanicOnError(err)
	//return pwUser.Entity, pwUser.RecCtx, nil
	return user, fs.RecCtxUser{}, nil
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
		return fs.RecCtxUser{}, errx.ErrxOf(err)
	}
	sql := `
	INSERT INTO users (username, email, password_hash, bio, image)
	VALUES ($1, $2, $3, $4, $5)
	`
	args := []any{user.Username, user.Email, user.PasswordHash, user.Bio, user.ImageLink}
	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fs.RecCtxUser{}, errx.ErrxOf(err)
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
