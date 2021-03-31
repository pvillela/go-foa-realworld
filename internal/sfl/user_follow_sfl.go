package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentAddSfl is the stereotype instance for the service flow that
// causes the current user start following a given other user.
type UserFollowSfl struct {
}

// UserFollowSflT is the function type instantiated by UserFollowSfl.
type UserFollowSflT = func(username string, followedUsername string) (*rpc.ProfileOut, error)

func (s UserFollowSfl) Make() UserFollowSflT {
	return func(username string, followedUsername string) (*rpc.ProfileOut, error) {
		panic("todo")
	}
}
