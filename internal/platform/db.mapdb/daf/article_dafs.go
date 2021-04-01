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
	return func(article model.Article) (fs.MdbArticle, error) {
		if _, err := s.getBySlug(article.Slug); err == nil {
			return fs.MdbArticle{}, fs.ErrDuplicateArticle
		}
		article.CreatedAt = time.Now()
		mdbArticle := fs.MdbArticle{
			Entity: article,
		}
		s.Store.Store(article.Slug, mdbArticle)
		return mdbArticle, nil
	}
}

func (s ArticleDafs) getBySlug(slug string) (fs.MdbArticle, error) {
	value, ok := s.Store.Load(slug)
	if !ok {
		return fs.MdbArticle{}, fs.ErrArticleNotFound
	}

	mdbArticle, ok := value.(fs.MdbArticle)
	if !ok {
		return fs.MdbArticle{}, errors.New("not an article stored at key")
	}

	return mdbArticle, nil
}

func (s ArticleDafs) MakeGetBySlug() fs.ArticleGetBySlugDafT {
	return s.getBySlug
}

func (s ArticleDafs) MakeUpdate() fs.ArticleUpdateDafT {
	return func(mdbArticle fs.MdbArticle) (fs.MdbArticle, error) {
		if _, err := s.getBySlug(mdbArticle.Slug); err != nil {
			return fs.MdbArticle{}, fs.ErrArticleNotFound
		}

		mdbArticle.UpdatedAt = time.Now()
		s.Store.Store(mdbArticle.Slug, mdbArticle)

		return mdbArticle, nil
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
}

func (s ArticleDafs) MakeGetRecentFiltered() fs.ArticleGetRecentFilteredDafT {
	return func(filters []model.ArticleFilter) ([]model.Article, error) {
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
}
