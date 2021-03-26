package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// GetProfileSflS contains the dependencies required for the construction of a
// GetProfileSfl. It represents the retrieval of a user profile.
type GetProfileSflS struct {
}

// GetProfileSfl is the type of a function that takes a string corresponding to a
// username and returns a model.ProfileOut.
type GetProfileSfl = func(username string) rpc.ProfileOut
