/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package model

import (
	"crypto/rand"
	"github.com/pvillela/go-foa-realworld/internal/arch/crypto"
	"sort"
	"time"
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
	FollowIDs      []string
	Favorites      []Article
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserUpdatableField int

const (
	UserName UserUpdatableField = iota
	UserEmail
	UserPassword
	UserBio
	UserImageLink
)

func (User) Create(
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
		FollowIDs:      nil,
		Favorites:      nil,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (s User) Update(fieldsToUpdate map[UserUpdatableField]interface{}) User {
	for k, v := range fieldsToUpdate {
		switch k {
		case UserName:
			s.Name = v.(string)
		case UserEmail:
			s.Email = v.(string)
		case UserPassword:
			password := v.(string)
			passwordHash := crypto.Hash(s.PasswordSalt, password)
			s.PasswordHash = passwordHash
		case UserBio:
			s.Bio = v.(*string)
		case UserImageLink:
			s.ImageLink = v.(string)
		}
	}
	s.UpdatedAt = time.Now()
	return s
}

func UpdateUser(user *User, opts ...func(fields *User)) {
	for _, v := range opts {
		v(user)
	}
}

func SetUserName(input string) func(fields *User) {
	return func(user *User) {
		user.Name = input
	}
}

func SetUserEmail(input string) func(fields *User) {
	return func(user *User) {
		user.Email = input
	}
}

func SetUserBio(input *string) func(fields *User) {
	return func(user *User) {
		user.Bio = input
	}
}

// give empty string to delete it
func SetUserImageLink(input string) func(fields *User) {
	return func(user *User) {
		user.ImageLink = input
	}
}

//func SetUserPassword(input *string) func(fields *Entity) {
//	return func(initial *Entity) {
//		if input != nil {
//			initial.Password = *input
//		}
//	}
//}

func (user User) Follows(userName string) bool {
	if user.FollowIDs == nil {
		return false
	}

	sort.Strings(user.FollowIDs)
	i := sort.SearchStrings(user.FollowIDs, userName)
	return i < len(user.FollowIDs) && user.FollowIDs[i] == userName
}

// UpdateFollowees will append or remove followee to current user according to follow param
func (s User) UpdateFollowees(followeeName string, follow bool) User {
	if follow {
		s.FollowIDs = append(s.FollowIDs, followeeName)
		return s
	}

	for i := 0; i < len(s.FollowIDs); i++ {
		if s.FollowIDs[i] == followeeName {
			s.FollowIDs = append(s.FollowIDs[:i], s.FollowIDs[i+1:]...) // memory leak ? https://github.com/golang/go/wiki/SliceTricks
		}
	}
	if len(s.FollowIDs) == 0 {
		s.FollowIDs = nil
	}
	return s
}
