/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	log "github.com/sirupsen/logrus"
	"strings"
)

// UserGetByNameDafI implements a stereotype instance of type
// UserGetByNameDafT.
var UserGetByNameDafI UserGetByNameDafT = func(
	ctx context.Context,
	username string,
) (model.User, RecCtxUser, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return model.User{}, RecCtxUser{}, err
	}
	return UserGetByNameExplicitTxDafI(ctx, tx, username)
}

// UserGetByNameExplicitTxDafI implements a stereotype instance of type
// UserGetByNameDafT with explicit passing of a pgx.Tx.
var UserGetByNameExplicitTxDafI UserGetByNameExplicitTxDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	username string,
) (model.User, RecCtxUser, error) {
	pwUser := PwUser{}
	err := dbpgx.ReadSingle(ctx, tx, "users", "username", username, &pwUser)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind == dbpgx.DbErrRecordNotFound {
			errX := util.SafeCast[errx.Errx](err)
			return model.User{}, RecCtxUser{}, errX.Customize(bf.ErrMsgUsernameNotFound, username)
		}
		return model.User{}, RecCtxUser{}, err
	}
	return pwUser.Entity, pwUser.RecCtx, nil
}

// UserGetByEmailDafI implements a stereotype instance of type
// UserGetByEmailDafT.
var UserGetByEmailDafI UserGetByEmailDafT = func(
	ctx context.Context,
	email string,
) (model.User, RecCtxUser, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return model.User{}, RecCtxUser{}, err
	}
	pwUser := PwUser{}
	err = dbpgx.ReadSingle(ctx, tx, "users", "email", strings.ToLower(email), &pwUser)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind == dbpgx.DbErrRecordNotFound {
			return model.User{}, RecCtxUser{},
				kind.Decorate(util.SafeCast[errx.Errx](err), bf.ErrMsgUserEmailNotFound, email)
		}
		return model.User{}, RecCtxUser{}, err
	}
	return pwUser.Entity, pwUser.RecCtx, nil
}

// UserCreateDafI implements a stereotype instance of type
// UserCreateDafT.
var UserCreateDafI UserCreateDafT = func(
	ctx context.Context,
	user *model.User,
) (RecCtxUser, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return RecCtxUser{}, err
	}
	return UserCreateExplicitTxDafI(ctx, tx, user)
}

// UserCreateExplicitTxDafI implements a stereotype instance of type
// UserCreateDafT.
var UserCreateExplicitTxDafI UserCreateExplicitTxDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	user *model.User,
) (RecCtxUser, error) {
	sql := `
	INSERT INTO users (username, email, password_hash, password_salt, bio, image)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at
	`
	user.Email = strings.ToLower(user.Email)
	args := []any{
		user.Username,
		user.Email,
		user.PasswordHash,
		user.PasswordSalt,
		user.Bio,
		user.ImageLink,
	}
	row := tx.QueryRow(ctx, sql, args...)
	var recCtx RecCtxUser
	err := row.Scan(&user.Id, &recCtx.CreatedAt, &recCtx.UpdatedAt)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind == dbpgx.DbErrUniqueViolation {
			return RecCtxUser{}, kind.Make(err, bf.ErrMsgUsernameDuplicate, user.Username)
		}
		return RecCtxUser{}, kind.Make(err, "")
	}

	return recCtx, nil
}

// UserUpdateDafI implements a stereotype instance of type
// UserUpdateDafT.
var UserUpdateDafI UserUpdateDafT = func(
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
	SET username = $1, email = $2, bio = $3, image = $4, password_hash = $5, password_salt = $6, 
		updated_at = NOW()
	WHERE id = $7 AND updated_at = $8
	RETURNING updated_at
	`
	user.Email = strings.ToLower(user.Email)
	args := []interface{}{
		user.Username,
		user.Email,
		user.Bio,
		user.ImageLink,
		user.PasswordHash,
		user.PasswordSalt,
		user.Id,
		recCtx.UpdatedAt,
	}
	log.Debug("UserUpdateDaf sql: ", sql)
	log.Debug("UserUpdateDaf args: ", args)

	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(&recCtx.UpdatedAt)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind == dbpgx.DbErrUniqueViolation {
			return RecCtxUser{}, kind.Make(err, bf.ErrMsgUsernameDuplicate, user.Username)
		}
		if kind == dbpgx.DbErrRecordNotFound {
			return RecCtxUser{}, kind.Make(err, bf.ErrMsgUsernameNotFound, user.Username)
		}
		return RecCtxUser{}, kind.Make(err, "")
	}

	return recCtx, nil
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
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind == dbpgx.DbErrRecordNotFound {
			return kind.Make(err, bf.ErrMsgUsernameNotFound, username)
		}
		return kind.Make(err, "")
	}

	return nil
}
