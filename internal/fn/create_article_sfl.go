package fn

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CreateArticleSflS contains the dependencies required for the construction of a
// CreateArticleSfl. It represents the creation of an article.
type CreateArticleSflS struct {
}

// CreateArticleSfl is the type of a function that takes an rpc.CreatArticleIn as input and
// returns a model.Article.
type CreateArticleSfl = func(articleIn rpc.CreateArticleIn) model.Article
