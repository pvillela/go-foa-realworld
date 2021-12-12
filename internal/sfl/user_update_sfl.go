/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserUpdateSflS is the stereotype instance for the service flow that
// It represents the action of registering a user.
type UserUpdateSflS struct {
	BeginTxn         func(context string) db.Txn
	UserGetByNameDaf fs.UserGetByNameDafT
	UserUpdateDaf    fs.UserUpdateDafT
}

// UserUpdateSflT is the type of a function that takes an rpc.UserUpdateIn as input
// and returns a model.User.
type UserUpdateSflT = func(username string, in rpc.UserUpdateIn) (rpc.UserOut, error)

func (s UserUpdateSflS) Make() UserUpdateSflT {
	userGenTokenBf := fs.UserGenTokenBfI
	return func(username string, in rpc.UserUpdateIn) (rpc.UserOut, error) {
		txn := s.BeginTxn("ArticleCreateSflS")
		defer txn.End()

		user, rc, err := s.UserGetByNameDaf(username)
		if err != nil {
			return rpc.UserOut{}, err
		}

		user = user.Update(in.User)

		_, err = s.UserUpdateDaf(user, rc, txn)
		if err != nil {
			return rpc.UserOut{}, err
		}

		token, err := userGenTokenBf(user)
		if err != nil {
			return rpc.UserOut{}, err
		}

		userOut := rpc.UserOut_FromModel(user, token)
		return userOut, err
	}
}
