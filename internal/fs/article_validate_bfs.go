package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/ft"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleValidateBeforeCreateBf struct{}

func (ArticleValidateBeforeCreateBf) invoke(article model.Article) error {
	return nil
}

func (s ArticleValidateBeforeCreateBf) Make() ft.ArticleValidateBeforeCreateBfT {
	return s.invoke
}

type ArticleValidateBeforeUpdateBf struct{}

func (ArticleValidateBeforeUpdateBf) invoke(article model.Article) error {
	return nil
}

func (s ArticleValidateBeforeUpdateBf) Make() ft.ArticleValidateBeforeUpdateBfT {
	return s.invoke
}
