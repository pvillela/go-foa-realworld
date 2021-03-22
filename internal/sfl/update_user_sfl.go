package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UpdateUserSflS contains the dependencies required for the construction of a
// UpdateUserSfl. It represents the action of updating user information.
type UpdateUserSflS struct {
}

// UpdateUserSfl is the type of a function that takes an rpc.UserUpdateIn as input
// and returns a model.User.
type UpdateUserSfl = func(userInfo rpc.UserUpdateIn) model.User
