/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserAuthenticateSflS is the stereotype instance for the service flow that
// authenticates a user.
type UserAuthenticateSflS struct {
	UserGetByEmailDaf  fs.UserGetByEmailDafT
	UserAuthenticateBf fs.UserAuthenticateBfT
}

// UserAuthenticateSflT is the function type instantiated by UserAuthenticateSflS.
type UserAuthenticateSflT = func(_username string, in rpc.UserAuthenticateIn) (rpc.UserOut, error)

func (s UserAuthenticateSflS) Make() UserAuthenticateSflT {
	userGenTokenBf := fs.UserGenTokenBfI
	return func(_ string, in rpc.UserAuthenticateIn) (rpc.UserOut, error) {
		var zero rpc.UserOut

		email := in.User.Email
		password := in.User.Password

		user, _, err := s.UserGetByEmailDaf(email)
		if err != nil {
			return zero, err
		}

		if !s.UserAuthenticateBf(user, password) {
			// I know the error info below is not secure but OK for now for debugging
			return zero, fs.ErrAuthenticationFailed.Make(nil, user.Name, password)
		}

		token, err := userGenTokenBf(user)
		if err != nil {
			return zero, err
		}

		userOut := rpc.UserOut_FromModel(user, token)
		return userOut, err
	}
}
