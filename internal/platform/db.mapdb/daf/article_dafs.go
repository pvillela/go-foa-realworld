/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package daf

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleDafs struct {
	ArticleDb mapdb.MapDb
}

func (s ArticleDafs) MakeCreate() fs.ArticleCreateDafT {
	return func(article model.Article, txn mapdb.Txn) (db.RecCtx, error) {
		_, _, err := s.getBySlug(article.Slug)
		if err == nil {
			return nil, fs.ErrDuplicateArticleSlug.Make(article.Slug)
		}
		if util.ErrKindOf(err) != mapdb.ErrRecordNotFound {
			return nil, err
		}

		pw := fs.PwArticle{nil, article}
		err = s.ArticleDb.Create(article.Uuid, pw, txn)
		if err != nil {
			return nil, err
		}
		return pw.RecCtx, nil
	}
}

func (s ArticleDafs) getBySlug(slug string) (model.Article, db.RecCtx, error) {
	var iVal interface{}
	var found bool
	s.ArticleDb.Range(func(key, value interface{}) bool {
		if key == slug {
			iVal = value
			found = true
			return false
		}
		return true
	})
	if !found {
		return model.Article{}, nil, fs.ErrArticleSlugNotFound.Make(slug)
	}
	pw, ok := iVal.(fs.PwArticle)
	if !ok {
		panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap article"))
	}

	return pw.Entity, pw.RecCtx, nil
}

func (s ArticleDafs) MakeGetBySlug() fs.ArticleGetBySlugDafT {
	return s.getBySlug
}

func (s ArticleDafs) MakeUpdate() fs.ArticleUpdateDafT {
	return func(article model.Article, recCtx db.RecCtx, txn mapdb.Txn) (db.RecCtx, error) {
		value, err := s.ArticleDb.Read(article.Uuid)
		if err != nil {
			return nil, err
		}
		pw, ok := value.(fs.PwArticle)
		if !ok {
			panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap article"))
		}
		if pw.Entity.Slug != article.Slug {
			return nil, fs.ErrDuplicateArticleSlug.Make(article.Slug)
		}

		pw.Entity = article
		err = s.ArticleDb.Update(article.Uuid, pw, txn)
		if err != nil {
			return nil, err
		}

		return pw.RecCtx, nil
	}
}

func (s ArticleDafs) MakeDelete() fs.ArticleDeleteDafT {
	return func(slug string, txn mapdb.Txn) error {
		article, _, err := s.getBySlug(slug)
		if err != nil {
			return err
		}

		err = s.ArticleDb.Delete(article.Uuid, txn)
		if err != nil {
			return err
		}

		return nil
	}
}

func (s ArticleDafs) MakeGetByAuthorsOrderedByMostRecentDaf() fs.ArticleGetByAuthorsOrderedByMostRecentDafT {
	return func(usernames []string) ([]model.Article, error) {
		var toReturn []model.Article

		s.ArticleDb.Range(func(key, value interface{}) bool {
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

		s.ArticleDb.Range(func(key, value interface{}) bool {
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
