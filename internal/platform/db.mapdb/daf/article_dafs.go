/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/newdaf"

	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

const (
	limitDefault  = 20
	offsetDefault = 0
)

type myMapDb struct {
	mapdb.MapDb
}

func pwArticleFromDb(value interface{}) fs.PwArticle {
	pw, ok := value.(fs.PwArticle)
	if !ok {
		panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap article"))
	}
	return pw
}

func articleFromDb(value interface{}) model.Article {
	return pwArticleFromDb(value).Entity
}

// ArticleCreateDafC is a function that constructs a stereotype instance of type
// fs.ArticleCreateDafT.
func ArticleCreateDafC(articleDb mapdb.MapDb) newdaf.ArticleCreateDafT {
	return func(article model.Article, txn db.Txn) (fs.RecCtxArticle, error) {
		_, _, err := myMapDb{articleDb}.getBySlug(article.Slug)
		if err == nil {
			return fs.RecCtxArticle{}, fs.ErrDuplicateArticleSlug.Make(nil, article.Slug)
		}

		pw := fs.PwArticle{fs.RecCtxArticle{}, article}
		err = articleDb.Create(article.Id, pw, txn)
		if errx.KindOf(err) == mapdb.ErrDuplicateKey {
			return fs.RecCtxArticle{}, fs.ErrDuplicateArticleUuid.Make(err, article.Id)
		}
		if err != nil {
			return fs.RecCtxArticle{}, err // this can only be a transaction error
		}

		return fs.RecCtxArticle{}, nil
	}
}

func (s myMapDb) getBySlug(slug string) (model.Article, fs.RecCtxArticle, error) {
	pred := func(_, value interface{}) bool {
		article := articleFromDb(value)
		if article.Slug == slug {
			return true
		}
		return false
	}

	value, found := s.FindFirst(pred)
	if !found {
		return model.Article{}, fs.RecCtxArticle{}, fs.ErrArticleSlugNotFound.Make(nil, slug)
	}
	pw := pwArticleFromDb(value)

	return pw.Entity, pw.RecCtx, nil
}

// ArticleGetBySlugDafC is a function that constructs a stereotype instance of type
// fs.ArticleGetBySlugDafT.
func ArticleGetBySlugDafC(articleDb mapdb.MapDb) newdaf.ArticleGetBySlugDafT {
	return myMapDb{articleDb}.getBySlug
}

// ArticleUpdateDafC is a function that constructs a stereotype instance of type
// fs.ArticleUpdateDafT.
func ArticleUpdateDafC(articleDb mapdb.MapDb) newdaf.ArticleUpdateDafT {
	return func(article model.Article, recCtx fs.RecCtxArticle, txn db.Txn) (fs.RecCtxArticle, error) {
		if artBySlug, _, err := (myMapDb{articleDb}.getBySlug(article.Slug)); err == nil && artBySlug.Id != article.Id {
			return fs.RecCtxArticle{}, fs.ErrDuplicateArticleSlug.Make(nil, article.Slug)
		}

		pw := fs.PwArticle{recCtx, article}
		err := articleDb.Update(article.Id, pw, txn)
		if errx.KindOf(err) == mapdb.ErrRecordNotFound {
			return fs.RecCtxArticle{}, fs.ErrArticleNotFound.Make(err, article.Id)
		}
		if err != nil {
			return fs.RecCtxArticle{}, err // this can only be a transaction error
		}

		return recCtx, nil
	}
}

// ArticleDeleteDafC is a function that constructs a stereotype instance of type
// fs.ArticleDeleteDafT.
func ArticleDeleteDafC(articleDb mapdb.MapDb) newdaf.ArticleDeleteDafT {
	return func(slug string, txn db.Txn) error {
		article, _, err := myMapDb{articleDb}.getBySlug(slug)
		if err != nil {
			return err
		}

		err = articleDb.Delete(article.Id, txn)
		if err != nil {
			return err // this can only be a transaction error because article was found above
		}

		return nil
	}
}

func selectAndOrderByMostRecent(
	articleDb mapdb.MapDb,
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

	selected := articleDb.FindAll(pred)
	less := func(i, j int) bool {
		return articleFromDb(selected[i]).CreatedAt.After(articleFromDb(selected[j]).CreatedAt)
	}
	util.Sort(selected, less)
	selected = util.SliceWindow(selected, limit, offset)

	res := make([]model.Article, limit)
	for i := range selected {
		res[i] = articleFromDb(selected[i])
	}

	return res
}

// ArticleGetRecentForAuthorsDafC is a function that constructs a stereotype instance of type
// fs.ArticleGetRecentForAuthorsDafT.
func ArticleGetRecentForAuthorsDafC(articleDb mapdb.MapDb) fs.ArticleGetRecentForAuthorsDafT {
	return func(usernames []string, pLimit, pOffset *int) ([]model.Article, error) {
		pred := func(_, value interface{}) bool {
			article := articleFromDb(value)
			for _, name := range usernames {
				if name == article.Author.Username {
					return true
				}
			}
			return false
		}
		res := selectAndOrderByMostRecent(articleDb, pred, pLimit, pOffset)
		return res, nil
	}
}

// ArticleGetRecentFilteredDafC is a function that constructs a stereotype instance of type
// fs.ArticleGetRecentFilteredDafT.
func ArticleGetRecentFilteredDafC(articleDb mapdb.MapDb) fs.ArticleGetRecentFilteredDafT {
	return func(in rpc.ArticlesListIn) ([]model.Article, error) {
		pred := func(_, value interface{}) bool {
			article := articleFromDb(value)

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

			if author := in.Author; author != nil && article.Author.Username != *author {
				return false
			}

			findFavoritedBy := func(favoritedBy string) bool {
				for _, user := range article.FavoritedBy {
					if user.Username == favoritedBy {
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

		res := selectAndOrderByMostRecent(articleDb, pred, in.Limit, in.Offset)
		return res, nil
	}
}
