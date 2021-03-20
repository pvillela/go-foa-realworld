package fn

import "github.com/pvillela/go-foa-realworld/internal/model"

// UnfollowUserSflS contains the dependencies required for the construction of a
// UnfollowUserSfl. It represents the action of having the current user stop following a given
// other user.
type UnfollowUserSflS struct {
}

// UnfollowUserSfl is the type of a function that takes the current username and a followed
// username and returns a model.Profile.
type UnfollowUserSfl = func(currentUsername string, followedUsername string) model.Profile
