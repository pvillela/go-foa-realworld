package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// ArticleGetAndCheckOwnerFl is the stereotype instance for the flow that
// checks if a given article's author's username matches a given username.
type ArticleGetAndCheckOwnerFl struct {
	ArticleGetBySlugDaf func(slug string) (*model.Article, error)
	CheckArticleOwnerBf func(article model.Article, username string) error
}

// ArticleGetAndCheckOwnerFlT is the function type instantiated by fs.ArticleGetAndCheckOwnerFl.
type ArticleGetAndCheckOwnerFlT = func(username, slug string) (*model.Article, error)

func (s ArticleGetAndCheckOwnerFl) invoke(slug string, username string) (*model.Article, error) {
	article, err := s.ArticleGetBySlugDaf(slug)
	if err != nil {
		return nil, err
	}

	if err := s.CheckArticleOwnerBf(*article, username); err != nil {
		return nil, err
	}

	return article, err
}

func (s ArticleGetAndCheckOwnerFl) Make() ArticleGetAndCheckOwnerFlT {
	return s.invoke
}
