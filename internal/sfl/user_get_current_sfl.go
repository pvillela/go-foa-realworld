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

// UserGetCurrentSfl is the stereotype instance for the service flow that
// returns the current user.
type UserGetCurrentSfl struct {
	UserGetByNameDaf fs.UserGetByNameDafT
	UserGenTokenBf   fs.UserGenTokenBfT
}

// UserGetCurrentSflT is the function type instantiated by UserGetCurrentSfl.
type UserGetCurrentSflT = func(username string) (rpc.UserOut, error)

func (s UserGetCurrentSfl) Make() UserGetCurrentSflT {
	return func(username string) (rpc.UserOut, error) {
		user, _, err := s.UserGetByNameDaf(username)
		if err != nil {
			return rpc.UserOut{}, err
		}

		token, err := s.UserGenTokenBf(user)
		if err != nil {
			return rpc.UserOut{}, err
		}

		userOut := rpc.UserOut_FromModel(user, token)
		return userOut, err
	}
}
