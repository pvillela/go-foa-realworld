package fn

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// AddCommentSflS contains the dependencies required for the construction of a
// AddCommentSfl. It represents the addition of a comment to an article.
type AddCommentSflS struct {
}

// AddCommentSfl is the type of a function that takes a slug and an rpc.AddCommentIn as inputs
// and returns a model.Comment.
type AddCommentSfl = func(slug string, commentIn rpc.AddCommentIn) model.Comment
