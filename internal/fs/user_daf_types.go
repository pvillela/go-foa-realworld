/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// PwUser is a wrapper of the model.User entity
// containing context information required for ersistence purposes.
type PwUser struct {
	db.RecCtx
	Entity model.User
}

type UserGetByNameDafT = func(userName string) (model.User, db.RecCtx, error)

type UserGetByEmailDafT = func(email string) (model.User, db.RecCtx, error)

type UserUpdateDafT = func(user model.User, recCtx db.RecCtx) (db.RecCtx, error)

type UserCreateDafT = func(user model.User) (db.RecCtx, error)
