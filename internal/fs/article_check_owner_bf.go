package fs

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleCheckOwnerBf struct{}

type ArticleCheckOwnerBfT = func(article model.Article, username string) error

func (ArticleCheckOwnerBf) Make() ArticleCheckOwnerBfT {
	return func(article model.Article, username string) error {
		if article.Author.Name != username {
			return errors.New("article not owned by user")
		}
		return nil
	}
}
