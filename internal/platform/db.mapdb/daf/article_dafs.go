/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package daf

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
)

type ArticleDafs struct {
	Store *sync.Map
}

func (s ArticleDafs) MakeCreate() fs.ArticleCreateDafT {
	return func(article model.Article) (db.RecCtx, error) {
		pw := fs.PwArticle{nil, article}
		_, loaded := s.Store.LoadOrStore(article.Slug, pw)
		if loaded {
			return nil, fs.ErrDuplicateArticle
		}
		return pw.RecCtx, nil
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
	return func(article model.Article, recCtx db.RecCtx) (db.RecCtx, error) {
		if _, _, err := s.getBySlug(article.Slug); err != nil {
			return nil, fs.ErrArticleNotFound
		}

		pw := fs.PwArticle{nil, article}
		s.Store.Store(article.Slug, pw)

		return pw.RecCtx, nil
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
