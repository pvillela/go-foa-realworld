/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// TODO: think about when to pass model.User by value vs by pointer. If User is cheap to copy
//  then we should always pass it by value unless we need to mutate the argument as in the case
//  of UserCreateDaf.

// UserGetByNameDaf implements a stereotype instance of type
// UserGetByNameDafT.
var UserGetByNameDaf UserGetByNameDafT = func(
	ctx context.Context,
	username string,
) (model.User, RecCtxUser, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return model.User{}, RecCtxUser{}, err
	}
	pwUser := PwUser{}
	err = dbpgx.ReadSingle(ctx, tx, "users", "username", username, &pwUser)
	if err != nil {
		return model.User{}, RecCtxUser{}, err
	}
	return pwUser.Entity, pwUser.RecCtx, nil
}

// UserGetByEmailDaf implements a stereotype instance of type
// UserGetByEmailDafT.
var UserGetByEmailDaf UserGetByEmailDafT = func(
	ctx context.Context,
	email string,
) (model.User, RecCtxUser, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return model.User{}, RecCtxUser{}, err
	}
	pwUser := PwUser{}
	err = dbpgx.ReadSingle(ctx, tx, "users", "email", email, &pwUser)
	if err != nil {
		return model.User{}, RecCtxUser{}, err
	}
	return pwUser.Entity, pwUser.RecCtx, nil
}

// UserCreateDaf implements a stereotype instance of type
// UserCreateDafT.
var UserCreateDaf UserCreateDafT = func(
	ctx context.Context,
	user *model.User,
) (RecCtxUser, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return RecCtxUser{}, errx.ErrxOf(err)
	}
	sql := `
	INSERT INTO users (username, email, password_hash, bio, image)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`
	args := []any{user.Username, user.Email, user.PasswordHash, user.Bio, user.ImageLink}
	row := tx.QueryRow(ctx, sql, args...)
	var recCtx RecCtxUser
	err = row.Scan(&user.Id, &recCtx.CreatedAt, &recCtx.UpdatedAt)
	return recCtx, errx.ErrxOf(err)
}

// UserUpdateDafC is the function that constructs a stereotype instance of type
// UserUpdateDafT.
var UserUpdateDaf UserUpdateDafT = func(
	ctx context.Context,
	user model.User,
	recCtx RecCtxUser,
) (RecCtxUser, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return RecCtxUser{}, errx.ErrxOf(err)
	}
	sql := `
	UPDATE users 
	SET username = $1, email = $2, bio = $3, image = $4, password_hash = $5, updated_at = NOW()
	WHERE id = $6 AND updated_at = $7
	RETURNING updated_at
	`
	args := []interface{}{
		user.Username,
		user.Email,
		user.Bio,
		user.ImageLink,
		user.PasswordHash,
		user.Id,
		recCtx.UpdatedAt,
	}
	fmt.Println("UserUpdateDaf sql:", sql)
	fmt.Println("UserUpdateDaf args:", args)
	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(&recCtx.UpdatedAt)
	return recCtx, errx.ErrxOf(err)
}

func userDeleteByUsernameDaf(
	ctx context.Context,
	username string,
) error {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return errx.ErrxOf(err)
	}
	sql := `
	DELETE FROM users
	WHERE username = $1
	`
	_, err = tx.Exec(ctx, sql, username)
	return errx.ErrxOf(err)
}
