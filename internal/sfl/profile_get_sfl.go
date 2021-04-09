package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ProfileGetSfl is the stereotype instance for the service flow that
// retrieves a user profile.
type ProfileGetSfl struct {
	UserGetByNameDaf fs.UserGetByNameDafT
}

// ProfileGetSflT is the function type instantiated by CommentsGetSfl.
type ProfileGetSflT = func(username, profileName string) (rpc.ProfileOut, error)

func (s ProfileGetSfl) Make() ProfileGetSflT {
	return func(username, profileName string) (rpc.ProfileOut, error) {
		var zero rpc.ProfileOut

		profileUser, _, err := s.UserGetByNameDaf(profileName)
		if err != nil {
			return zero, err
		}

		var follows bool
		if username != "" {
			user, _, err := s.UserGetByNameDaf(username)
			if err != nil {
				return zero, err
			}
			follows = user.Follows(profileName)
		}

		profileOut := rpc.ProfileOut{}.FromModel(profileUser, follows)

		return profileOut, nil
	}
}
