package fn

import "github.com/pvillela/go-foa-realworld/internal/model"

type GetArticleAndCheckOwnerFl struct {
	ArticleGetBySlugDaf         func(slug string) (*model.Article, error)
	CheckArticleUserOwnershipBf func(article model.Article, username string) error
}

func (s GetArticleAndCheckOwnerFl) Invoke(slug string, username string) (*model.Article, error) {
	article, err := s.ArticleGetBySlugDaf(slug)
	if err != nil {
		return nil, err
	}

	if err := s.CheckArticleUserOwnershipBf(*article, username); err != nil {
		return nil, err
	}

	return article, err
}
