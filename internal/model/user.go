/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package model

import (
	"sort"

	"github.com/pvillela/go-foa-realworld/internal/arch/crypto"
)

// User represents a user account in the system
type User struct {
	Id           uint    `json:"-"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	PasswordHash string  `json:"-"`
	Bio          *string `json:"bio,omitempty"`
	ImageLink    string  `json:"image,omitempty"`
	Followees    []*User `json:"-"`
	Followers    []*User `json:"-"`
	// Below added to daf.RecCtx
	//CreatedAt      time.Time `json:"-"`
	//UpdatedAt      time.Time `json:"-"`
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
	passwordHash := crypto.ArgonPasswordHash(password)

	return User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Bio:          nil,
		ImageLink:    "",
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

func (user User) Follows(userName string) bool {
	if user.Followees == nil {
		return false
	}

	sort.Strings(user.Followees)
	i := sort.SearchStrings(user.Followees, userName)
	return i < len(user.Followees) && user.Followees[i] == userName
}

// UpdateFollowees appends or removes followee to current user according to follow param
func (s User) UpdateFollowees(followeeName string, follow bool) User {
	if follow {
		s.Followees = append(s.Followees, followeeName)
		return s
	}

	for i := 0; i < len(s.Followees); i++ {
		if s.Followees[i] == followeeName {
			s.Followees = append(s.Followees[:i], s.Followees[i+1:]...) // TODO: memory leak ? https://github.com/golang/go/wiki/SliceTricks
			break
		}
	}
	if len(s.Followees) == 0 {
		s.Followees = nil
	}
	return s
}
