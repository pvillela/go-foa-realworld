package sfl

import "github.com/pvillela/go-foa-realworld/internal/model"

// GetTagsSflS contains the dependencies required for the construction of a
// GetTagsSfl. It represents the retrieval of all tags.
type GetTagsSflS struct {
}

// GetTagsSfl is the type of a function that takes no inputs and
// and returns a model.Tags.
type GetTagsSfl = func() model.Tags
