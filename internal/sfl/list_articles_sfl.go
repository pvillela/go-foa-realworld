package sfl

import "github.com/pvillela/go-foa-realworld/internal/model"

// ListArticlesSflS contains the dependencies required for the construction of a
// ListArticlesSfl. It represents a query for articles based on a set of query parameters.
type ListArticlesSflS struct {
}

// ListArticlesSfl is the type of a function that takes a number of optional query
// parameters and returns a model.Articles.
type ListArticlesSfl = func(
	tag string,
	author string,
	favorited string,
	limit int,
	offset int,
) model.Articles
