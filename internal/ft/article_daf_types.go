package ft

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleCreateDafT = func(article model.Article) (*model.Article, error)

type ArticleGetBySlugDafT = func(slug string) (*model.Article, error)

type ArticleUpdateDafT = func(article model.Article) (*model.Article, error)

type ArticleDeleteDafT = func(slug string) error

type ArticleGetByAuthorsOrderedByMostRecentDafT = func(usernames []string) ([]model.Article, error)

type ArticleGetRecentFilteredDafT = func(filters []model.ArticleFilter) ([]model.Article, error)
