package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// FollowUserSflS contains the dependencies required for the construction of a
// FollowUserSfl. It represents the action of  having the current user start following a given
// other user.
type FollowUserSflS struct {
}

// FollowUserSfl is the type of a function that takes the current username and a followed
// username and returns a model.ProfileOut.
type FollowUserSfl = func(currentUsername string, followedUsername string) rpc.ProfileOut
