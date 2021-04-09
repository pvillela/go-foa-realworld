package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentAddSfl is the stereotype instance for the service flow that
// causes the current user start following a given other user.
type UserFollowSfl struct {
	UserFollowFl fs.UserFollowFlT
}

// UserFollowSflT is the function type instantiated by UserFollowSfl.
type UserFollowSflT = func(username string, followedUsername string) (rpc.ProfileOut, error)

func (s UserFollowSfl) Make() UserFollowSflT {
	return func(username string, followedUsername string) (rpc.ProfileOut, error) {
		var zero rpc.ProfileOut
		user, _, err := s.UserFollowFl(username, followedUsername, true)
		if err != nil {
			return zero, err
		}
		profileOut := rpc.ProfileOut{}.FromModel(user, true)
		return profileOut, err
	}
}
