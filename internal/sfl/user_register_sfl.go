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

// UserRegisterSflS is the stereotype instance for the service flow that
// It represents the action of registering a user.
type UserRegisterSflS struct {
	BeginTxn      func(context string) db.Txn
	UserCreateDaf fs.UserCreateDafT
}

// UserRegisterSflT is the type of a function that takes an rpc.UserRegisterIn as input
// and returns a model.User.
type UserRegisterSflT = func(in rpc.UserRegisterIn) (rpc.UserOut, error)

func (s UserRegisterSflS) Make() UserRegisterSflT {
	userGenTokenBf := fs.UserGenTokenBfI
	return func(in rpc.UserRegisterIn) (rpc.UserOut, error) {
		txn := s.BeginTxn("ArticleCreateSflS")
		defer txn.End()

		user := in.ToUser()

		_, err := s.UserCreateDaf(user, txn)
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
