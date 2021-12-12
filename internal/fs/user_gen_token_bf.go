/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type UserGenTokenBfT = func(user model.User) (string, error)

var UserGenTokenBfI UserGenTokenBfT = func(user model.User) (string, error) {
	return UserGenTokenSup(user)
}
