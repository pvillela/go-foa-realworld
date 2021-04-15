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
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

type ArticleDafs struct {
	ArticleDb mapdb.MapDb
}

const (
	limitDefault  = 20
	offsetDefault = 0
)

func (s ArticleDafs) MakeCreate() fs.ArticleCreateDafT {
	return func(article model.Article, txn db.Txn) (db.RecCtx, error) {
		_, _, err := s.getBySlug(article.Slug)
		if err == nil {
			return nil, fs.ErrDuplicateArticleSlug.Make(nil, article.Slug)
		}

		pw := fs.PwArticle{nil, article}
		err = s.ArticleDb.Create(article.Uuid, pw, txn)
		if util.ErrKindOf(err) == mapdb.ErrDuplicateKey {
			return nil, fs.ErrDuplicateArticleUuid.Make(err, article.Uuid)
		}
		if err != nil {
			return nil, err // this can only be a transaction error
		}

		return nil, nil
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
		return model.Article{}, nil, fs.ErrArticleSlugNotFound.Make(nil, slug)
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
	return func(article model.Article, recCtx db.RecCtx, txn db.Txn) (db.RecCtx, error) {
		if artBySlug, _, err := s.getBySlug(article.Slug); err == nil && artBySlug.Uuid != article.Uuid {
			return nil, fs.ErrDuplicateArticleSlug.Make(nil, article.Slug)
		}

		pw := fs.PwArticle{recCtx, article}
		err := s.ArticleDb.Update(article.Uuid, pw, txn)
		if util.ErrKindOf(err) == mapdb.ErrRecordNotFound {
			return nil, fs.ErrArticleNotFound.Make(err, article.Uuid)
		}
		if err != nil {
			return nil, err // this can only be a transaction error
		}

		return recCtx, nil
	}
}

func (s ArticleDafs) MakeDelete() fs.ArticleDeleteDafT {
	return func(slug string, txn db.Txn) error {
		article, _, err := s.getBySlug(slug)
		if err != nil {
			return err
		}

		err = s.ArticleDb.Delete(article.Uuid, txn)
		if err != nil {
			return err // this can only be a transaction error because article was found above
		}

		return nil
	}
}

func (s ArticleDafs) selectAndOrderByMostRecent(
	pred func(key, value interface{}) bool,
	pLimit, pOffset *int,
) []model.Article {
	limit := limitDefault
	if pLimit != nil {
		limit = *pLimit
	}
	offset := offsetDefault
	if pOffset != nil {
		offset = *pOffset
	}

	selected := s.ArticleDb.FindAll(pred)
	less := func(i, j int) bool {
		return selected[i].(fs.PwArticle).Entity.CreatedAt.After(selected[i].(fs.PwArticle).Entity.CreatedAt)
	}
	util.Sort(selected, less)
	selected = util.SliceWindow(selected, limit, offset)

	res := make([]model.Article, limit)
	for i := range selected {
		res[i] = selected[i].(fs.PwArticle).Entity
	}

	return res
}

func (s ArticleDafs) MakeGetRecentForAuthorsDaf() fs.ArticleGetRecentForAuthorsDafT {
	return func(usernames []string, pLimit, pOffset *int) ([]model.Article, error) {
		pred := func(_, value interface{}) bool {
			article := value.(fs.PwArticle).Entity
			for _, name := range usernames {
				if name == article.Author.Name {
					return true
				}
			}
			return false
		}
		res := s.selectAndOrderByMostRecent(pred, pLimit, pOffset)
		return res, nil
	}
}

func (s ArticleDafs) MakeGetRecentFiltered() fs.ArticleGetRecentFilteredDafT {
	return func(in rpc.ArticlesListIn) ([]model.Article, error) {
		pred := func(_, value interface{}) bool {
			article := value.(fs.PwArticle).Entity

			findTag := func(tag string) bool {
				for _, t := range article.TagList {
					if t == tag {
						return true
					}
				}
				return false
			}
			if tag := in.Tag; tag != nil && !findTag(*tag) {
				return false
			}

			if author := in.Author; author != nil && article.Author.Name != *author {
				return false
			}

			findFavoritedBy := func(favoritedBy string) bool {
				for _, user := range article.FavoritedBy {
					if user.Name == favoritedBy {
						return true
					}
				}
				return false
			}
			if favorited := in.Favorited; favorited != nil && findFavoritedBy(*favorited) {
				return false
			}

			return true
		}

		res := s.selectAndOrderByMostRecent(pred, in.Limit, in.Offset)
		return res, nil
	}
}
