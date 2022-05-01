/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/newdaf"

	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

func pwUserFromDb(value interface{}) newdaf.PwUser {
	pw, ok := value.(newdaf.PwUser)
	if !ok {
		panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap user"))
	}
	return pw
}

func userFromDb(value interface{}) model.User {
	return pwUserFromDb(value).Entity
}

func getByName(
	userDb mapdb.MapDb,
	username string,
) (model.User, newdaf.RecCtxUser, error) {
	value, err := userDb.Read(username)
	if err != nil {
		return model.User{}, newdaf.RecCtxUser{}, fs.ErrUserNameNotFound.Make(err, username)
	}
	pw := pwUserFromDb(value)
	return pw.Entity, pw.RecCtx, nil
}

// UserGetByNameDafC is the function that constructs a stereotype instance of type
// fs.UserGetByNameDafT.
func UserGetByNameDafC(
	userDb mapdb.MapDb,
) newdaf.UserGetByNameDafT {
	return func(userName string) (model.User, newdaf.RecCtxUser, error) {
		return getByName(userDb, userName)
	}
}

func getByEmail(userDb mapdb.MapDb, email string) (model.User, newdaf.RecCtxUser, error) {
	pred := func(_, value interface{}) bool {
		if userFromDb(value).Email == email {
			return true
		}
		return false
	}

	value, found := userDb.FindFirst(pred)
	if !found {
		return model.User{}, newdaf.RecCtxUser{}, fs.ErrUserEmailNotFound.Make(nil, email)
	}
	pw := pwUserFromDb(value)
	return pw.Entity, pw.RecCtx, nil
}

// UserGetByEmailDafC is the function that constructs a stereotype instance of type
// fs.UserGetByEmailDafT.
func UserGetByEmailDafC(
	userDb mapdb.MapDb,
) newdaf.UserGetByEmailDafT {
	return func(email string) (model.User, newdaf.RecCtxUser, error) {
		return getByEmail(userDb, email)
	}
}

// UserCreateDafC is the function that constructs a stereotype instance of type
// fs.UserCreateDafT.
func UserCreateDafC(
	userDb mapdb.MapDb,
) newdaf.UserCreateDafT {
	return func(user model.User, txn db.Txn) (newdaf.RecCtxUser, error) {
		if _, _, err := getByEmail(userDb, user.Email); err == nil {
			return newdaf.RecCtxUser{}, fs.ErrDuplicateUserEmail.Make(nil, user.Email)
		}

		pwUser := newdaf.PwUser{RecCtx: newdaf.RecCtxUser{}, Entity: user}
		err := userDb.Create(user.Username, pwUser, txn)
		if errx.KindOf(err) == mapdb.ErrDuplicateKey {
			return newdaf.RecCtxUser{}, fs.ErrDuplicateUserName.Make(err, user.Username)
		}
		if err != nil {
			return newdaf.RecCtxUser{}, err
		}

		return pwUser.RecCtx, nil
	}
}

// UserUpdateDafC is the function that constructs a stereotype instance of type
// fs.UserUpdateDafT.
func UserUpdateDafC(
	userDb mapdb.MapDb,
) newdaf.UserUpdateDafT {
	return func(user model.User, recCtx newdaf.RecCtxUser, txn db.Txn) (newdaf.RecCtxUser, error) {
		if userByEmail, _, err := getByEmail(userDb, user.Email); err == nil && userByEmail.Username != user.Username {
			return newdaf.RecCtxUser{}, fs.ErrDuplicateUserEmail.Make(nil, user.Email)
		}

		pw := newdaf.PwUser{RecCtx: recCtx, Entity: user}
		err := userDb.Update(user.Username, pw, txn)
		if errx.KindOf(err) == mapdb.ErrRecordNotFound {
			return newdaf.RecCtxUser{}, fs.ErrUserNameNotFound.Make(err, user.Username)
		}
		if err != nil {
			return newdaf.RecCtxUser{}, err // this can only be a transaction error
		}

		return recCtx, nil
	}
}
