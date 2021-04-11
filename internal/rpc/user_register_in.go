/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import (
	"crypto/rand"
	"github.com/pvillela/go-foa-realworld/internal/arch/jwt"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"time"
)

type UserRegisterIn struct {
	User struct {
		Username string
		Email    string
		Password string
	}
}

func (s UserRegisterIn) ToUser() model.User {
	c := 16
	salt := make([]byte, c)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err) // should never happen
	}
	passwordHash := jwt.Hash(salt, s.User.Password)
	user := model.User{
		Name:           s.User.Username,
		Email:          s.User.Email,
		IsTempPassword: false,
		PasswordHash:   passwordHash,
		PasswordSalt:   salt,
		Bio:            nil,
		ImageLink:      "",
		FollowIDs:      nil,
		Favorites:      nil,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return user
}
