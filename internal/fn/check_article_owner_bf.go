package fn

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type CheckArticleOwnerBf struct{}

func (CheckArticleOwnerBf) Invoke(article model.Article, username string) error {
	if article.Author.Name != username {
		return errors.New("article not owned by user")
	}
	return nil
}
