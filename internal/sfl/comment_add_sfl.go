package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentAddSflS contains the dependencies required for the construction of a
// CommentAddSfl. It represents the addition of a comment to an article.
type CommentAddSflS struct {
}

// CommentAddSfl is the type of a function that takes a slug and an rpc.CommentAddIn as inputs
// and returns a model.Comment.
type CommentAddSfl = func(slug string, commentIn rpc.CommentAddIn) model.Comment
