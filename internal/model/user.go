/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package model

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/crypto"
	"strings"
)

// User represents a user account in the system
type User struct {
	Id           uint    `json:"-"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	PasswordHash string  `json:"-"`
	Bio          *string `json:"bio,omitempty"`
	ImageLink    string  `json:"image,omitempty" db:"image"`
	//Followees    []*User `json:"-"`
	//Followers    []*User `json:"-"`
	// Below added to daf.RecCtx
	//CreatedAt time.Time `json:"-"`
	//UpdatedAt time.Time `json:"-"`
}

type Profile struct {
	Username  string  `json:"username"`
	Bio       *string `json:"bio"`
	Image     *string `json:"image"`
	Following bool    `json:"following" db:"-"`
}

type UserPatch struct {
	Username  *string
	Email     *string
	Password  *string
	Bio       *string
	ImageLink *string
}

func User_Create(
	username string,
	email string,
	password string,
) User {
	passwordHash := crypto.ArgonPasswordHash(password)

	return User{
		Username:     username,
		Email:        strings.ToLower(email),
		PasswordHash: passwordHash,
		Bio:          nil,
		ImageLink:    "",
	}
}

func (s User) Update(v UserPatch) User {
	if v.Username != nil {
		s.Username = *v.Username
	}
	if v.Email != nil {
		s.Email = strings.ToLower(*v.Email)
	}
	if v.Password != nil {
		password := *v.Password
		passwordHash := crypto.ArgonPasswordHash(password)
		s.PasswordHash = passwordHash
	}
	if v.Bio != nil {
		s.Bio = v.Bio
	}
	if v.ImageLink != nil {
		s.ImageLink = *v.ImageLink
	}

	return s
}

func ProfileFromUser(user *User, follows bool) Profile {
	s := Profile{}
	if user.Bio != nil {
		s.Bio = user.Bio
	} else {
		empty := ""
		s.Bio = &empty
	}

	if user.ImageLink != "" {
		s.Image = &user.ImageLink
	}

	s.Username = user.Username
	s.Following = follows

	return s
}
