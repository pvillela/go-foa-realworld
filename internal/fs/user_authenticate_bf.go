/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/crypto"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type UserAuthenticateBf struct{}

type UserAuthenticateBfT = func(user model.User, password string) bool

func (UserAuthenticateBf) Make() UserAuthenticateBfT {
	return func(user model.User, password string) bool {
		if crypto.Hash(user.PasswordSalt, password) != user.PasswordHash {
			return false
		}
		return true
	}
}
