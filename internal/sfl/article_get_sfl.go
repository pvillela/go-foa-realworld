package sfl

import "github.com/pvillela/go-foa-realworld/internal/model"

// GetArticleSflS contains the dependencies required for the construction of a
// GetArticleSfl. It represents the retrieval of a single article.
type GetArticleSflS struct {
}

// GetArticleSfl is the type of a function that takes a string corresponding to an
// article slug and returns a model.Article.
type GetArticleSfl = func(slug string) model.Article
