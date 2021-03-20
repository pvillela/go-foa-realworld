package fn

import "github.com/pvillela/go-foa-realworld/internal/model"

// GetCurrentUserSflS contains the dependencies required for the construction of a
// GetCurrentUserSfl. It represents the action of having the current user start following a
// given other user.
type GetCurrentUserSflS struct {
}

// GetCurrentUsersSfl is the type of a function that takes a JWT token corresponding to the
// current user and returns that user.
type GetCurrentUserSfl = func(token string) model.User
