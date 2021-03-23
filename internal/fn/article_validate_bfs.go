package fn

import "github.com/pvillela/go-foa-realworld/internal/model"

type ArticleValidateBeforeCreateBf struct{}

func (ArticleValidateBeforeCreateBf) Invoke(article model.Article) error {
	return nil
}

type ArticleValidateBeforeUpdateBf struct{}

func (ArticleValidateBeforeUpdateBf) Invoke(article model.Article) error {
	return nil
}
