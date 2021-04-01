package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type MdbArticle struct {
	db.RecCtx
	Entity model.Article
}

type ArticleCreateDafT = func(article model.Article) (MdbArticle, error)

type ArticleGetBySlugDafT = func(slug string) (MdbArticle, error)

type ArticleUpdateDafT = func(mdbArticle MdbArticle) (MdbArticle, error)

type ArticleDeleteDafT = func(slug string) error

type ArticleGetByAuthorsOrderedByMostRecentDafT = func(usernames []string) ([]MdbArticle, error)

type ArticleGetRecentFilteredDafT = func(filters []model.ArticleFilter) ([]MdbArticle, error)
