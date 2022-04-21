/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fs

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// PwUser is a type alias
type PwUser = db.Pw[model.User]

// RecCtxUser is a type alias
type RecCtxUser = db.RecCtx[model.User]

// UserGetByNameDafT is the type of the stereotype instance for the DAF that
// retrieves a user by username.
type UserGetByNameDafT = func(reqCtx context.Context, userName string) (model.User, RecCtxUser, error)

// UserGetByEmailDafT is the type of the stereotype instance for the DAF that
// retrieves a user by email address.
type UserGetByEmailDafT = func(reqCtx context.Context, email string) (model.User, RecCtxUser, error)

// UserCreateDafT is the type of the stereotype instance for the DAF that
// creates a user.
type UserCreateDafT = func(reqCtx context.Context, user model.User, txn db.Txn) (RecCtxUser, error)

// UserUpdateDafT is the type of the stereotype instance for the DAF that
// updates a user.
type UserUpdateDafT = func(user model.User, recCtx RecCtxUser, txn db.Txn) (RecCtxUser, error)
