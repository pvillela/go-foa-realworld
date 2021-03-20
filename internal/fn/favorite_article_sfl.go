package fn

import "github.com/pvillela/go-foa-realworld/internal/model"

// FavoriteArticleSflS contains the dependencies required for the construction of a
// FavoriteArticleSflS. It represents the designation of an article as a favorite.
type FavoriteArticleSflS struct {
}

// FavoriteArticleSfl is the type of a function that takes a slug as input and
// returns a model.Article.
type FavoriteArticleSfl = func(slug string) model.Article
