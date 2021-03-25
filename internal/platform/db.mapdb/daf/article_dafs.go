package daf

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/ft"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
	"time"
)

type ArticleDafs struct {
	Store *sync.Map
}

func (s ArticleDafs) create(article model.Article) (*model.Article, error) {
	if _, err := s.getBySlug(article.Slug); err == nil {
		return nil, fs.ErrDuplicateArticle
	}
	article.CreatedAt = time.Now()
	s.Store.Store(article.Slug, article)
	return &article, nil
}

func (s ArticleDafs) MakeCreate() ft.ArticleCreateDafT {
	return s.create
}

func (s ArticleDafs) getBySlug(slug string) (*model.Article, error) {
	value, ok := s.Store.Load(slug)
	if !ok {
		return nil, fs.ErrArticleNotFound
	}

	article, ok := value.(model.Article)
	if !ok {
		return nil, errors.New("not an article stored at key")
	}

	return &article, nil
}

func (s ArticleDafs) MakeGetBySlug() ft.ArticleGetBySlugDafT {
	return s.getBySlug
}

func (s ArticleDafs) update(article model.Article) (*model.Article, error) {
	if _, err := s.getBySlug(article.Slug); err != nil {
		return nil, fs.ErrArticleNotFound
	}

	article.UpdatedAt = time.Now()
	s.Store.Store(article.Slug, article)

	return &article, nil
}

func (s ArticleDafs) MakeUpdate() ft.ArticleUpdateDafT {
	return s.update
}

func (s ArticleDafs) delete(slug string) error {
	s.Store.Delete(slug)

	return nil
}

func (s ArticleDafs) MakeDelete() ft.ArticleDeleteDafT {
	return s.delete
}

func (s ArticleDafs) getByAuthorsOrderedByMostRecent(usernames []string) ([]model.Article, error) {
	var toReturn []model.Article

	s.Store.Range(func(key, value interface{}) bool {
		article, ok := value.(model.Article)
		if !ok {
			return true // log this but continue
		}
		for _, username := range usernames {
			if article.Author.Name == username {
				toReturn = append(toReturn, article)
			}
		}
		return true
	})

	return toReturn, nil
}

func (s ArticleDafs) MakeGetByAuthorsOrderedByMostRecentDaf() ft.ArticleGetByAuthorsOrderedByMostRecentDafT {
	return s.getByAuthorsOrderedByMostRecent
}

func (s ArticleDafs) getRecentFiltered(filters []model.ArticleFilter) ([]model.Article, error) {
	var recentArticles []model.Article

	s.Store.Range(func(key, value interface{}) bool {
		article, ok := value.(model.Article)
		if !ok {
			// not an article (shouldn't happen) -> skip
			return true // log this but continue
		}

		for _, funcToApply := range filters {
			if !funcToApply(article) { // "AND filter" : if one of the filter is at false, skip the article
				return true
			}
		}

		recentArticles = append(recentArticles, article)
		return true
	})

	return recentArticles, nil
}

func (s ArticleDafs) MakeGetRecentFiltered() ft.ArticleGetRecentFilteredDafT {
	return s.getRecentFiltered
}
