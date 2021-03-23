package fn

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
	"time"
)

type ArticleDafs struct {
	Store *sync.Map
}

func (s ArticleDafs) Create(article model.Article) (*model.Article, error) {
	if _, err := s.GetBySlug(article.Slug); err == nil {
		return nil, ErrDuplicateArticle
	}
	article.CreatedAt = time.Now()
	s.Store.Store(article.Slug, article)
	return &article, nil
}

func (s ArticleDafs) GetBySlug(slug string) (*model.Article, error) {
	value, ok := s.Store.Load(slug)
	if !ok {
		return nil, ErrArticleNotFound
	}

	article, ok := value.(model.Article)
	if !ok {
		return nil, errors.New("not an article stored at key")
	}

	return &article, nil
}

func (s ArticleDafs) Update(article model.Article) (*model.Article, error) {
	if _, err := s.GetBySlug(article.Slug); err != nil {
		return nil, ErrArticleNotFound
	}

	article.UpdatedAt = time.Now()
	s.Store.Store(article.Slug, article)

	return &article, nil
}

func (s ArticleDafs) Delete(slug string) error {
	s.Store.Delete(slug)

	return nil
}
