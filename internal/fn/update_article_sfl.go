package fn

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UpdateArticleSflS contains the dependencies required for the construction of a
// UpdateArticleSfl. It represents the updating of an article.
type UpdateArticleSflS struct {
}

// UpdateArticleSfl is the type of a function that takes an rpc.UpdateArticleIn as input and
// returns a model.Article.
type UpdateArticleSfl = func(articleIn rpc.UpdateArticleIn) model.Article
