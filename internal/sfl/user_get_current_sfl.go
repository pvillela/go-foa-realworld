package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserGetCurrentSfl is the stereotype instance for the service flow that
// returns the current user.
type UserGetCurrentSfl struct {
}

// UserGetCurrentSflT is the function type instantiated by UserGetCurrentSfl.
type UserGetCurrentSflT = func(username string) rpc.UserOut
