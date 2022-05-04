/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

func pwUserFromDb(value interface{}) daf.PwUser {
	pw, ok := value.(daf.PwUser)
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
) (model.User, daf.RecCtxUser, error) {
	value, err := userDb.Read(username)
	if err != nil {
		return model.User{}, daf.RecCtxUser{}, bf.ErrUserNameNotFound.Make(err, username)
	}
	pw := pwUserFromDb(value)
	return pw.Entity, pw.RecCtx, nil
}

// UserGetByNameDafC is the function that constructs a stereotype instance of type
// bf.UserGetByNameDafT.
func UserGetByNameDafC(
	userDb mapdb.MapDb,
) daf.UserGetByNameDafT {
	return func(userName string) (model.User, daf.RecCtxUser, error) {
		return getByName(userDb, userName)
	}
}

func getByEmail(userDb mapdb.MapDb, email string) (model.User, daf.RecCtxUser, error) {
	pred := func(_, value interface{}) bool {
		if userFromDb(value).Email == email {
			return true
		}
		return false
	}

	value, found := userDb.FindFirst(pred)
	if !found {
		return model.User{}, daf.RecCtxUser{}, bf.ErrUserEmailNotFound.Make(nil, email)
	}
	pw := pwUserFromDb(value)
	return pw.Entity, pw.RecCtx, nil
}

// UserGetByEmailDafC is the function that constructs a stereotype instance of type
// bf.UserGetByEmailDafT.
func UserGetByEmailDafC(
	userDb mapdb.MapDb,
) daf.UserGetByEmailDafT {
	return func(email string) (model.User, daf.RecCtxUser, error) {
		return getByEmail(userDb, email)
	}
}

// UserCreateDafC is the function that constructs a stereotype instance of type
// bf.UserCreateDafT.
func UserCreateDafC(
	userDb mapdb.MapDb,
) daf.UserCreateDafT {
	return func(user model.User, txn db.Txn) (daf.RecCtxUser, error) {
		if _, _, err := getByEmail(userDb, user.Email); err == nil {
			return daf.RecCtxUser{}, bf.ErrDuplicateUserEmail.Make(nil, user.Email)
		}

		pwUser := daf.PwUser{RecCtx: daf.RecCtxUser{}, Entity: user}
		err := userDb.Create(user.Username, pwUser, txn)
		if errx.KindOf(err) == mapdb.ErrDuplicateKey {
			return daf.RecCtxUser{}, bf.ErrDuplicateUserName.Make(err, user.Username)
		}
		if err != nil {
			return daf.RecCtxUser{}, err
		}

		return pwUser.RecCtx, nil
	}
}

// UserUpdateDafC is the function that constructs a stereotype instance of type
// bf.UserUpdateDafT.
func UserUpdateDafC(
	userDb mapdb.MapDb,
) daf.UserUpdateDafT {
	return func(user model.User, recCtx daf.RecCtxUser, txn db.Txn) (daf.RecCtxUser, error) {
		if userByEmail, _, err := getByEmail(userDb, user.Email); err == nil && userByEmail.Username != user.Username {
			return daf.RecCtxUser{}, bf.ErrDuplicateUserEmail.Make(nil, user.Email)
		}

		pw := daf.PwUser{RecCtx: recCtx, Entity: user}
		err := userDb.Update(user.Username, pw, txn)
		if errx.KindOf(err) == mapdb.ErrRecordNotFound {
			return daf.RecCtxUser{}, bf.ErrUserNameNotFound.Make(err, user.Username)
		}
		if err != nil {
			return daf.RecCtxUser{}, err // this can only be a transaction error
		}

		return recCtx, nil
	}
}
