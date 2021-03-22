package sfl

import "github.com/pvillela/go-foa-realworld/internal/model"

// FeedArticlesSflS contains the dependencies required for the construction of a
// FeedArticlesSfl. It represents a query for all articles authored by other users followed by
// the current user, with optional limit and offset pagination parameters.
type FeedArticlesSflS struct {
}

// FeedArticlesSfl is the type of a function that takes optional pagination parameters
// and returns a model.Articles.
type FeedArticlesSfl = func(
	limit int,
	offset int,
) model.Articles
