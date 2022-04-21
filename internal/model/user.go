/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package model

import (
	"sort"
	"time"

	"github.com/pvillela/go-foa-realworld/internal/arch/crypto"
)

// User represents a user account in the system
type User struct {
	Id             uint   `json:"-"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	IsTempPassword bool
	PasswordHash   string    `json:"-"`
	Bio            *string   `json:"bio,omitempty"`
	ImageLink      string    `json:"image,omitempty"`
	Following      []*User   `json:"-"`
	Followers      []*User   `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

type UserUpdateSrc struct {
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
	now := time.Now()
	c := 16
	passwordHash := crypto.BcryptPasswordHash(password)

	return User{
		Username:       username,
		Email:          email,
		IsTempPassword: false,
		PasswordHash:   passwordHash,
		Bio:            nil,
		ImageLink:      "",
		Following:      nil,
		//Favorites:      nil,
		CreatedAt:    now,
		UpdatedAt:    now,
		NumFollowers: 0,
	}
}

func (s User) Update(v UserUpdateSrc) User {
	if v.Username != nil {
		s.Username = *v.Username
	}
	if v.Email != nil {
		s.Email = *v.Email
	}
	if v.Password != nil {
		password := *v.Password
		passwordHash := crypto.BcryptPasswordHash(s.PasswordSalt, password)
		s.PasswordHash = passwordHash
	}
	if v.Bio != nil {
		s.Bio = v.Bio
	}
	if v.ImageLink != nil {
		s.ImageLink = *v.ImageLink
	}

	s.UpdatedAt = time.Now()
	return s
}

func (user User) Follows(userName string) bool {
	if user.Following == nil {
		return false
	}

	sort.Strings(user.Following)
	i := sort.SearchStrings(user.Following, userName)
	return i < len(user.Following) && user.Following[i] == userName
}

// UpdateFollowees appends or removes followee to current user according to follow param
func (s User) UpdateFollowees(followeeName string, follow bool) User {
	if follow {
		s.Following = append(s.Following, followeeName)
		return s
	}

	for i := 0; i < len(s.Following); i++ {
		if s.Following[i] == followeeName {
			s.Following = append(s.Following[:i], s.Following[i+1:]...) // TODO: memory leak ? https://github.com/golang/go/wiki/SliceTricks
			break
		}
	}
	if len(s.Following) == 0 {
		s.Following = nil
	}
	return s
}
