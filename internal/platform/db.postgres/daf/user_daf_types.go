/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// RecCtxUser is a type alias
type RecCtxUser = dbpgx.RecCtx[model.User]

// PwUser is a type alias
type PwUser = db.Pw[model.User, RecCtxUser]

// UserGetByNameDafT is the type of the stereotype instance for the DAF that
// retrieves a user by username.
type UserGetByNameDafT = func(ctx context.Context, userName string) (model.User, RecCtxUser, error)

// UserGetByNameExplicitTxDafT is the type of the stereotype instance for the DAF that
// retrieves a user by username taking an explicit pgx.Tx.
type UserGetByNameExplicitTxDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	username string,
) (model.User, RecCtxUser, error)

// UserGetByEmailDafT is the type of the stereotype instance for the DAF that
// retrieves a user by email address.
type UserGetByEmailDafT = func(ctx context.Context, email string) (model.User, RecCtxUser, error)

// UserCreateDafT is the type of the stereotype instance for the DAF that
// creates a user.
type UserCreateDafT = func(ctx context.Context, user *model.User) (RecCtxUser, error)

// UserCreateExplicitTxDafT is the type of the stereotype instance for the DAF that
// creates a user taking an explicit pgx.Tx.
type UserCreateExplicitTxDafT = func(ctx context.Context, tx pgx.Tx, user *model.User) (RecCtxUser, error)

// UserUpdateDafT is the type of the stereotype instance for the DAF that
// updates a user.
type UserUpdateDafT = func(ctx context.Context, user model.User, recCtx RecCtxUser) (RecCtxUser, error)
