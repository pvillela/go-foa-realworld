package fs

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/ft"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleCheckOwnerBf struct{}

func (ArticleCheckOwnerBf) invoke(article model.Article, username string) error {
	if article.Author.Name != username {
		return errors.New("article not owned by user")
	}
	return nil
}

func (s ArticleCheckOwnerBf) Make() ft.ArticleCheckOwnerBfT {
	return s.invoke
}
