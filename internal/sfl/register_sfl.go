package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// RegisterSflS contains the dependencies required for the construction of a
// RegisterSfl. It represents the action of registering a user.
type RegisterSflS struct {
}

// RegisterSfl is the type of a function that takes an rpc.UserRegisterIn as input
// and returns a model.User.
type RegisterSfl = func(in rpc.UserRegisterIn) model.User
