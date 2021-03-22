package sfl

import "github.com/pvillela/go-foa-realworld/internal/model"

// GetProfileSflS contains the dependencies required for the construction of a
// GetProfileSfl. It represents the retrieval of a user profile.
type GetProfileSflS struct {
}

// GetProfileSfl is the type of a function that takes a string corresponding to a
// username and returns a model.Profile.
type GetProfileSfl = func(username string) model.Profile
