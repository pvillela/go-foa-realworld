package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
)

// ArticleFavoriteSfl is the stereotype instance for the service flow that
// deletes an article.
type ArticleDeleteSfl struct {
	ArticleGetAndCheckOwnerFl fs.ArticleGetAndCheckOwnerFlT
	ArticleDeleteDaf          fs.ArticleDeleteDafT
}

// ArticleDeleteSflT is the function type instantiated by ArticleDeleteSfl.
type ArticleDeleteSflT = func(username, slug string) error

func (s ArticleDeleteSfl) Make() ArticleDeleteSflT {
	return func(username string, slug string) error {
		_, err := s.ArticleGetAndCheckOwnerFl(username, slug)
		if err != nil {
			return err
		}

		return s.ArticleDeleteDaf(slug)
	}
}
