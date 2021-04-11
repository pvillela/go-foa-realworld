/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserRegisterSfl is the stereotype instance for the service flow that
// It represents the action of registering a user.
type UserRegisterSfl struct {
	UserCreateDaf  fs.UserCreateDafT
	UserGenTokenBf fs.UserGenTokenBfT
}

// UserRegisterSflT is the type of a function that takes an rpc.UserRegisterIn as input
// and returns a model.User.
type UserRegisterSflT = func(in rpc.UserRegisterIn) (rpc.UserOut, error)

func (s UserRegisterSfl) Make() UserRegisterSflT {
	return func(in rpc.UserRegisterIn) (rpc.UserOut, error) {
		user := in.ToUser()

		_, err := s.UserCreateDaf(user)
		if err != nil {
			return rpc.UserOut{}, err
		}

		token, err := s.UserGenTokenBf(user)
		if err != nil {
			return rpc.UserOut{}, err
		}

		userOut := rpc.UserOut{}.FromModel(user, token)
		return userOut, err
	}
}
