/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package model

import (
	"crypto/rand"
	"sort"
	"time"

	"github.com/pvillela/go-foa-realworld/internal/arch/crypto"
)

// User represents a user account in the system
type User struct {
	Name           string
	Email          string
	IsTempPassword bool
	PasswordHash   string
	PasswordSalt   []byte
	Bio            *string
	ImageLink      string
	FollowedNames  []string // usernames
	NumFollowers   int
	CreatedAt      time.Time
	UpdatedAt      time.Time
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
	salt := make([]byte, c)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err) // should never happen
	}
	passwordHash := crypto.Hash(salt, password)

	return User{
		Name:           username,
		Email:          email,
		IsTempPassword: false,
		PasswordHash:   passwordHash,
		PasswordSalt:   salt,
		Bio:            nil,
		ImageLink:      "",
		FollowedNames:  nil,
		//Favorites:      nil,
		CreatedAt:    now,
		UpdatedAt:    now,
		NumFollowers: 0,
	}
}

func (s User) Update(v UserUpdateSrc) User {
	if v.Username != nil {
		s.Name = *v.Username
	}
	if v.Email != nil {
		s.Email = *v.Email
	}
	if v.Password != nil {
		password := *v.Password
		passwordHash := crypto.Hash(s.PasswordSalt, password)
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
	if user.FollowedNames == nil {
		return false
	}

	sort.Strings(user.FollowedNames)
	i := sort.SearchStrings(user.FollowedNames, userName)
	return i < len(user.FollowedNames) && user.FollowedNames[i] == userName
}

// UpdateFollowees appends or removes followee to current user according to follow param
func (s User) UpdateFollowees(followeeName string, follow bool) User {
	if follow {
		s.FollowedNames = append(s.FollowedNames, followeeName)
		return s
	}

	for i := 0; i < len(s.FollowedNames); i++ {
		if s.FollowedNames[i] == followeeName {
			s.FollowedNames = append(s.FollowedNames[:i], s.FollowedNames[i+1:]...) // TODO: memory leak ? https://github.com/golang/go/wiki/SliceTricks
			break
		}
	}
	if len(s.FollowedNames) == 0 {
		s.FollowedNames = nil
	}
	return s
}
