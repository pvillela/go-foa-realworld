package daf

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
	"time"
)

type ArticleDafs struct {
	Store *sync.Map
}

func (s ArticleDafs) MakeCreate() fs.ArticleCreateDafT {
	return func(article model.Article) (model.Article, db.RecCtx, error) {
		if _, _, err := s.getBySlug(article.Slug); err == nil {
			return model.Article{}, nil, fs.ErrDuplicateArticle
		}
		article.CreatedAt = time.Now()
		pw := fs.PwArticle{nil, article}
		s.Store.Store(article.Slug, pw)
		return pw.Entity, pw.RecCtx, nil
	}
}

func (s ArticleDafs) getBySlug(slug string) (model.Article, db.RecCtx, error) {
	value, ok := s.Store.Load(slug)
	if !ok {
		return model.Article{}, nil, fs.ErrArticleNotFound
	}

	pw, ok := value.(fs.PwArticle)
	if !ok {
		panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap article"))
	}

	return pw.Entity, pw.RecCtx, nil
}

func (s ArticleDafs) MakeGetBySlug() fs.ArticleGetBySlugDafT {
	return s.getBySlug
}

func (s ArticleDafs) MakeUpdate() fs.ArticleUpdateDafT {
	return func(article model.Article, recCtx db.RecCtx) (model.Article, db.RecCtx, error) {
		if _, _, err := s.getBySlug(article.Slug); err != nil {
			return model.Article{}, nil, fs.ErrArticleNotFound
		}

		article.UpdatedAt = time.Now()
		pw := fs.PwArticle{nil, article}
		s.Store.Store(article.Slug, pw)

		return pw.Entity, pw.RecCtx, nil
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
			pw, ok := value.(fs.PwArticle)
			if !ok {
				panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap article"))
			}
			for _, username := range usernames {
				if pw.Entity.Author.Name == username {
					toReturn = append(toReturn, pw.Entity)
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
			pw, ok := value.(fs.PwArticle)
			if !ok {
				panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap article"))
			}

			for _, funcToApply := range filters {
				if !funcToApply(pw.Entity) { // "AND filter" : if one of the filter is at false, skip the pw
					return true
				}
			}

			toReturn = append(toReturn, pw.Entity)
			return true
		})

		return toReturn, nil
	}
}
