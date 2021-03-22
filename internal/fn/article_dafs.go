package fn

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
)

type ArticleDafPre struct {
	Store *sync.Map
}

type ArticleDafS struct {
	ArticleDafPre
	Other struct{}
}

type ArticleCreateDaf = func(article model.Article) (*model.Article, error)

func (pre ArticleDafPre) Prep() ArticleDafS {
	return ArticleDafS{
		ArticleDafPre: pre,
	}
}

func (s ArticleDafS) Create(article model.Article) (*model.Article, error) {
	if _, err := s.GetBySlug(article.Slug); err == nil {
		return nil, ErrDuplicateArticle
	}
	s.Store.Store(article.Slug, article)
	return &article, nil
}

type ArticleGetBySlugDaf = func(string) (*model.Article, error)

func (s ArticleDafS) GetBySlug(slug string) (*model.Article, error) {
	return &model.Article{}, nil
}
