package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// PwArticle is a wrapper of the model.User entity
// containing context information required for persistence purposes.
type PwArticle interface {
	Entity() model.Article
	Updated(model.Article) PwArticle
}

type ArticleCreateDafT = func(article model.Article) (PwArticle, error)

type ArticleGetBySlugDafT = func(slug string) (PwArticle, error)

type ArticleUpdateDafT = func(pwArticle PwArticle) (PwArticle, error)

type ArticleDeleteDafT = func(slug string) error

type ArticleGetByAuthorsOrderedByMostRecentDafT = func(usernames []string) ([]model.Article, error)

type ArticleGetRecentFilteredDafT = func(filters []model.ArticleFilter) ([]model.Article, error)
