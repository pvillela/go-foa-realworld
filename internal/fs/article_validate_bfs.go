package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleValidateBeforeCreateBf struct{}

type ArticleValidateBeforeCreateBfT = func(article model.Article) error

func (ArticleValidateBeforeCreateBf) invoke(article model.Article) error {
	return nil
}

func (s ArticleValidateBeforeCreateBf) Make() ArticleValidateBeforeCreateBfT {
	return s.invoke
}

type ArticleValidateBeforeUpdateBf struct{}

type ArticleValidateBeforeUpdateBfT = func(article model.Article) error

func (ArticleValidateBeforeUpdateBf) invoke(article model.Article) error {
	return nil
}

func (s ArticleValidateBeforeUpdateBf) Make() ArticleValidateBeforeUpdateBfT {
	return s.invoke
}
