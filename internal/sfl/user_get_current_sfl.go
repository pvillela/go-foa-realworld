/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/arch"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleCreateSflT is the type of the stereotype instance for the service flow that
// returns the current user.
type UserGetCurrentSflT = func(username string, _ arch.Unit) (rpc.UserOut, error)

// UserGetCurrentSflC is the function that constructs a stereotype instance of type
// UserGetCurrentSflT.
func UserGetCurrentSflC(
	userGetByNameDaf fs.UserGetByNameDafT,
) UserGetCurrentSflT {
	userGenTokenBf := fs.UserGenTokenBfI
	return func(username string, _ arch.Unit) (rpc.UserOut, error) {
		user, _, err := userGetByNameDaf(username)
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
