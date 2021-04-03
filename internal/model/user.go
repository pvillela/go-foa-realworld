package model

import (
	"sort"
	"time"
)

// User represents a user account in the system
type User struct {
	Name         string
	Email        string
	TempPassword string
	PasswordHash string
	PasswordSalt []byte
	Bio          *string
	ImageLink    string
	FollowIDs    []string
	Favorites    []Article
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserUpdatableProperty int

const (
	UserEmail UserUpdatableProperty = iota
	UserName
	UserBio
	UserImageLink
	UserPassword
)

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
