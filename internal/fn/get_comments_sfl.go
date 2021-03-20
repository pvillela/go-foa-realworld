package fn

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// GetCommentsSflS contains the dependencies required for the construction of a
// GetCommentsSfl. It represents the retrieval of comments of an article.
type GetCommentsSflS struct {
}

// AddCommentSfl is the type of a function that takes a slug as input
// and returns a model.Comments.
type GetCommentsSfl = func(slug string) model.Comments
