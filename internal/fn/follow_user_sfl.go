package fn

import "github.com/pvillela/go-foa-realworld/internal/model"

// FollowUserSflS contains the dependencies required for the construction of a
// FollowUserSfl. It represents the action of  having the current user start following a given
// other user.
type FollowUserSflS struct {
}

// FollowUserSfl is the type of a function that takes the current username and a followed
// username and returns a model.Profile.
type FollowUserSfl = func(currentUsername string, followedUsername string) model.Profile
