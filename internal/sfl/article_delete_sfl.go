package sfl

import "github.com/pvillela/go-foa-realworld/internal/model"

// ArticleDeleteSflS contains the dependencies required for the construction of a
// ArticleDeleteSflT. It represents the deletion of an article.
type ArticleDeleteSfl struct {
	GetArticleAndCheckOwnerFl func(slug string, username string) (*model.Article, error)
	ArticleDeleteDaf          func(slug string) error
}

// ArticleDeleteSflT is the type of a function that takes a slug as input and
// returns noting.
type ArticleDeleteSflT = func(slug string)

func (s ArticleDeleteSfl) core(username string, slug string) error {
	_, err := s.GetArticleAndCheckOwnerFl(username, slug)
	if err != nil {
		return err
	}

	return s.ArticleDeleteDaf(slug)
}
