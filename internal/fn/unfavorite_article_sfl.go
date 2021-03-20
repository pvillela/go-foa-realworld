package fn

import "github.com/pvillela/go-foa-realworld/internal/model"

// UnfavoriteArticleSflS contains the dependencies required for the construction of a
// UnfavoriteArticleSfl. It represents the removal of the designation of an article as a favorite.
type UnfavoriteArticleSflS struct {
}

// UnfavoriteArticleSfl is the type of a function that takes a slug as input and
// returns a model.Article.
type UnfavoriteArticleSfl = func(slug string) model.Article
