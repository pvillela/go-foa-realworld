package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// AuthenticateSflS contains the dependencies required for the construction of a
// AuthenticateSfl. It represents the action of authenticating a user.
type AuthenticateSflS struct {
}

// AuthenticateSfl is the type of a function that takes an rpc.AuthenticateIn as input
// and returns a model.User.
type AuthenticateSfl = func(auth rpc.AuthenticateIn) model.User
