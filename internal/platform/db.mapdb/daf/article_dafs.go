package daf

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
	"time"
)

type ArticleDafs struct {
	Store *sync.Map
}

func (s ArticleDafs) MakeCreate() fs.ArticleCreateDafT {
	return func(article model.Article) (fs.PwArticle, error) {
		if _, err := s.getBySlug(article.Slug); err == nil {
			return fs.PwArticle{}, fs.ErrDuplicateArticle
		}
		article.CreatedAt = time.Now()
		pwArticle := fs.PwArticle{nil, article}
		s.Store.Store(article.Slug, pwArticle)
		return pwArticle, nil
	}
}

func (s ArticleDafs) getBySlug(slug string) (fs.PwArticle, error) {
	value, ok := s.Store.Load(slug)
	if !ok {
		return fs.PwArticle{}, fs.ErrArticleNotFound
	}

	pwArticle, ok := value.(fs.PwArticle)
	if !ok {
		return fs.PwArticle{}, errors.New("not an article stored at key")
	}

	return pwArticle, nil
}

func (s ArticleDafs) MakeGetBySlug() fs.ArticleGetBySlugDafT {
	return s.getBySlug
}

func (s ArticleDafs) MakeUpdate() fs.ArticleUpdateDafT {
	return func(pwArticle fs.PwArticle) (fs.PwArticle, error) {
		article := &pwArticle.Entity
		if _, err := s.getBySlug(article.Slug); err != nil {
			return fs.PwArticle{}, fs.ErrArticleNotFound
		}

		article.UpdatedAt = time.Now()
		s.Store.Store(article.Slug, pwArticle)

		return pwArticle, nil
	}
}

func (s ArticleDafs) MakeDelete() fs.ArticleDeleteDafT {
	return func(slug string) error {
		s.Store.Delete(slug)

		return nil
	}
}

func (s ArticleDafs) MakeGetByAuthorsOrderedByMostRecentDaf() fs.ArticleGetByAuthorsOrderedByMostRecentDafT {
	return func(usernames []string) ([]model.Article, error) {
		var toReturn []model.Article

		s.Store.Range(func(key, value interface{}) bool {
			pwArticle, ok := value.(fs.PwArticle)
			if !ok {
				return true // log this but continue
			}
			for _, username := range usernames {
				if pwArticle.Entity.Author.Name == username {
					toReturn = append(toReturn, pwArticle.Entity)
				}
			}
			return true
		})

		return toReturn, nil
	}
}

func (s ArticleDafs) MakeGetRecentFiltered() fs.ArticleGetRecentFilteredDafT {
	return func(filters []model.ArticleFilter) ([]model.Article, error) {
		var toReturn []model.Article

		s.Store.Range(func(key, value interface{}) bool {
			pwArticle, ok := value.(fs.PwArticle)
			if !ok {
				// not an pwArticle (shouldn't happen) -> skip
				return true // log this but continue
			}

			for _, funcToApply := range filters {
				if !funcToApply(pwArticle.Entity) { // "AND filter" : if one of the filter is at false, skip the pwArticle
					return true
				}
			}

			toReturn = append(toReturn, pwArticle.Entity)
			return true
		})

		return toReturn, nil
	}
}
