package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserAuthenticateSflS contains the dependencies required for the construction of a
// UserAuthenticateSfl. It represents the action of authenticating a user.
type UserAuthenticateSflS struct {
}

// UserAuthenticateSfl is the type of a function that takes an rpc.UserAuthenticateIn as input
// and returns a model.User.
type UserAuthenticateSfl = func(auth rpc.UserAuthenticateIn) model.User
