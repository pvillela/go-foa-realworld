package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleValidateBeforeCreateBf struct{}

type ArticleValidateBeforeCreateBfT = func(article model.Article) error

func (s ArticleValidateBeforeCreateBf) Make() ArticleValidateBeforeCreateBfT {
	return func(article model.Article) error {
		return nil
	}
}

type ArticleValidateBeforeUpdateBf struct{}

type ArticleValidateBeforeUpdateBfT = func(article model.Article) error

func (s ArticleValidateBeforeUpdateBf) Make() ArticleValidateBeforeUpdateBfT {
	return func(article model.Article) error {
		return nil
	}
}
