/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import "github.com/pvillela/go-foa-realworld/experimental/model"

type UserOut struct {
	User UserOut0 `json:"user"`
}

type UserOut0 = struct {
	Email    string  `json:"email"`
	Token    string  `json:"token"`
	Username string  `json:"username"`
	Bio      *string `json:"bio"`
	Image    string  `json:"image"`
}

func UserOut_FromModel(user model.User, token string) UserOut {
	return UserOut{
		User: UserOut0{
			Email:    user.Email,
			Token:    token,
			Username: user.Username,
			Bio:      user.Bio,
			Image:    user.ImageLink,
		},
	}
}
