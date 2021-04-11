/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import "github.com/pvillela/go-foa-realworld/internal/model"

type UserOut struct {
	User userOut0
}

type userOut0 struct {
	Email    string
	Token    string
	Username string
	Bio      *string
	Image    *string
}

func (s UserOut) FromModel(user model.User, token string) UserOut {
	s.User = userOut0{
		Email:    user.Email,
		Token:    token,
		Username: user.Name,
		Bio:      user.Bio,
	}
	if link := user.ImageLink; link != "" {
		s.User.Image = &link
	}
	return s
}
