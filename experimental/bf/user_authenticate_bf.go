/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package bf

import (
	"github.com/pvillela/go-foa-realworld/experimental/arch/crypto"
	"github.com/pvillela/go-foa-realworld/experimental/model"
)

type UserAuthenticateBfT = func(user model.User, password string) bool

var UserAuthenticateBf UserAuthenticateBfT = func(
	user model.User,
	password string,
) bool {
	return crypto.ArgonPasswordCheck(password, user.PasswordHash)
}
