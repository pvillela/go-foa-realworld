/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserAuthenticateSflT is the type of the stereotype instance for the service flow that
// authenticates a user.
type UserAuthenticateSflT = func(_ context.Context, in rpc.UserAuthenticateIn) (rpc.UserOut, error)

// UserAuthenticateSflC is the function that constructs a stereotype instance of type
// UserAuthenticateSflT.
func UserAuthenticateSflC(
	userGetByEmailDaf fs.UserGetByEmailDafT,
) UserAuthenticateSflT {
	userGenTokenBf := fs.UserGenTokenBfI
	userAuthenticateBf := fs.UserAuthenticateBfI
	return func(_ context.Context, in rpc.UserAuthenticateIn) (rpc.UserOut, error) {
		var zero rpc.UserOut

		email := in.User.Email
		password := in.User.Password

		user, _, err := userGetByEmailDaf(email)
		if err != nil {
			return zero, err
		}

		if !userAuthenticateBf(user, password) {
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
