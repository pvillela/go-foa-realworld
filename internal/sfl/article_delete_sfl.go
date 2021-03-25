package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
)

// ArticleFavoriteSfl is the stereotype instance for the service flow that
// deletes an article.
type ArticleDeleteSfl struct {
	GetArticleAndCheckOwnerFl fs.ArticleGetAndCheckOwnerFlT
	ArticleDeleteDaf          func(slug string) error
}

// ArticleDeleteSflT is the function type instantiated by ArticleDeleteSfl.
type ArticleDeleteSflT = func(username, slug string) error

func (s ArticleDeleteSfl) invoke(username string, slug string) error {
	_, err := s.GetArticleAndCheckOwnerFl(username, slug)
	if err != nil {
		return err
	}

	return s.ArticleDeleteDaf(slug)
}

func (s ArticleDeleteSfl) Make() ArticleDeleteSflT {
	return s.invoke
}
