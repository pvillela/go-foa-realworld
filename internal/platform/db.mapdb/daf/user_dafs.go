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

type UserDafsS struct {
	UserDb mapdb.MapDb
}

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

func (s UserDafsS) getByName(username string) (model.User, db.RecCtx, error) {
	value, err := s.UserDb.Read(username)
	if err != nil {
		return model.User{}, nil, fs.ErrUserNameNotFound.Make(err, username)
	}
	pw := pwUserFromDb(value)
	return pw.Entity, pw.RecCtx, nil
}

func (s UserDafsS) MakeGetByName() fs.UserGetByNameDafT {
	return s.getByName
}

func (s UserDafsS) getByEmail(email string) (model.User, db.RecCtx, error) {
	pred := func(_, value interface{}) bool {
		if userFromDb(value).Email == email {
			return true
		}
		return false
	}

	value, found := s.UserDb.FindFirst(pred)
	if !found {
		return model.User{}, nil, fs.ErrUserEmailNotFound.Make(nil, email)
	}
	pw := pwUserFromDb(value)
	return pw.Entity, pw.RecCtx, nil
}

func (s UserDafsS) MakeGetByEmail() fs.UserGetByEmailDafT {
	return func(email string) (model.User, db.RecCtx, error) {
		return s.getByEmail(email)
	}
}

func (s UserDafsS) MakeCreate() fs.UserCreateDafT {
	return func(user model.User, txn db.Txn) (db.RecCtx, error) {
		if _, _, err := s.getByEmail(user.Email); err == nil {
			return nil, fs.ErrDuplicateUserEmail.Make(nil, user.Email)
		}

		pwUser := fs.PwUser{RecCtx: nil, Entity: user}
		err := s.UserDb.Create(user.Name, pwUser, txn)
		if errx.KindOf(err) == mapdb.ErrDuplicateKey {
			return nil, fs.ErrDuplicateUserName.Make(err, user.Name)
		}
		if err != nil {
			return nil, err
		}

		return pwUser.RecCtx, nil
	}
}

func (s UserDafsS) MakeUpdate() fs.UserUpdateDafT {
	return func(user model.User, recCtx db.RecCtx, txn db.Txn) (db.RecCtx, error) {
		if userByEmail, _, err := s.getByEmail(user.Email); err == nil && userByEmail.Name != user.Name {
			return nil, fs.ErrDuplicateUserEmail.Make(nil, user.Email)
		}

		pw := fs.PwUser{RecCtx: recCtx, Entity: user}
		err := s.UserDb.Update(user.Name, pw, txn)
		if errx.KindOf(err) == mapdb.ErrRecordNotFound {
			return nil, fs.ErrUserNameNotFound.Make(err, user.Name)
		}
		if err != nil {
			return nil, err // this can only be a transaction error
		}

		return recCtx, nil
	}
}
