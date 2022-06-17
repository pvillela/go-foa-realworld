/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type UserRegisterIn struct {
	User UserRegisterIn0
}

type UserRegisterIn0 struct {
	Username string
	Email    string
	Password string
}

func (s UserRegisterIn) ToUser() model.User {
	return model.User_Create(s.User.Username, s.User.Email, s.User.Password)
}
