package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// GetCommentsSflS contains the dependencies required for the construction of a
// GetCommentsSfl. It represents the retrieval of comments of an article.
type GetCommentsSflS struct {
}

// CommentAddSflT is the type of a function that takes a slug as input
// and returns a model.Comments.
type GetCommentsSfl = func(slug string) rpc.CommentsOut
