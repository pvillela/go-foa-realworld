/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"fmt"

	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

func pwUserFromDb(value interface{}) fs.PwUser {
	pw, ok := value.(fs.PwUser)
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
) (model.User, fs.RecCtxUser, error) {
	value, err := userDb.Read(username)
	if err != nil {
		return model.User{}, fs.RecCtxUser{}, fs.ErrUserNameNotFound.Make(err, username)
	}
	pw := pwUserFromDb(value)
	return pw.Entity, pw.RecCtx, nil
}

// UserGetByNameDafC is the function that constructs a stereotype instance of type
// fs.UserGetByNameDafT.
func UserGetByNameDafT(
	userDb mapdb.MapDb,
) fs.UserGetByNameDafT {
	return func(userName string) (model.User, fs.RecCtxUser, error) {
		return getByName(userDb, userName)
	}
}

func getByEmail(userDb mapdb.MapDb, email string) (model.User, fs.RecCtxUser, error) {
	pred := func(_, value interface{}) bool {
		if userFromDb(value).Email == email {
			return true
		}
		return false
	}

	value, found := userDb.FindFirst(pred)
	if !found {
		return model.User{}, fs.RecCtxUser{}, fs.ErrUserEmailNotFound.Make(nil, email)
	}
	pw := pwUserFromDb(value)
	return pw.Entity, pw.RecCtx, nil
}

// UserGetByEmailDafC is the function that constructs a stereotype instance of type
// fs.UserGetByEmailDafT.
func UserGetByEmailDafC(
	userDb mapdb.MapDb,
) fs.UserGetByEmailDafT {
	return func(email string) (model.User, fs.RecCtxUser, error) {
		return getByEmail(userDb, email)
	}
}

// UserCreateDafC is the function that constructs a stereotype instance of type
// fs.UserCreateDafT.
func UserCreateDafC(
	userDb mapdb.MapDb,
) fs.UserCreateDafT {
	return func(user model.User, txn db.Txn) (fs.RecCtxUser, error) {
		if _, _, err := getByEmail(userDb, user.Email); err == nil {
			return fs.RecCtxUser{}, fs.ErrDuplicateUserEmail.Make(nil, user.Email)
		}

		pwUser := fs.PwUser{RecCtx: fs.RecCtxUser{}, Entity: user}
		err := userDb.Create(user.Name, pwUser, txn)
		if errx.KindOf(err) == mapdb.ErrDuplicateKey {
			return fs.RecCtxUser{}, fs.ErrDuplicateUserName.Make(err, user.Name)
		}
		if err != nil {
			return fs.RecCtxUser{}, err
		}

		return pwUser.RecCtx, nil
	}
}

// UserUpdateDafC is the function that constructs a stereotype instance of type
// fs.UserUpdateDafT.
func UserUpdateDafC(
	userDb mapdb.MapDb,
) fs.UserUpdateDafT {
	return func(user model.User, recCtx fs.RecCtxUser, txn db.Txn) (fs.RecCtxUser, error) {
		if userByEmail, _, err := getByEmail(userDb, user.Email); err == nil && userByEmail.Name != user.Name {
			return fs.RecCtxUser{}, fs.ErrDuplicateUserEmail.Make(nil, user.Email)
		}

		pw := fs.PwUser{RecCtx: recCtx, Entity: user}
		err := userDb.Update(user.Name, pw, txn)
		if errx.KindOf(err) == mapdb.ErrRecordNotFound {
			return fs.RecCtxUser{}, fs.ErrUserNameNotFound.Make(err, user.Name)
		}
		if err != nil {
			return fs.RecCtxUser{}, err // this can only be a transaction error
		}

		return recCtx, nil
	}
}
