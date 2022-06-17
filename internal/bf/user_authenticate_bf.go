/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package bf

import (
	"github.com/alexedwards/argon2id"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type UserAuthenticateBfT = func(user model.User, password string) bool

var UserAuthenticateBfI UserAuthenticateBfT = func(
	user model.User,
	password string,
) bool {
	ok, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	errx.PanicOnError(err)
	return ok
}
