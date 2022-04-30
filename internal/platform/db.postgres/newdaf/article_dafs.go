/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package newdaf

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

const (
	limitDefault  = 20
	offsetDefault = 0
)

// ArticleCreateDaf implements a stereotype instance of type
// fs.ArticleCreateDafT.
var ArticleCreateDaf fs.ArticleCreateDafT = func(
	ctx context.Context,
	article *model.Article,
) error {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return errx.ErrxOf(err)
	}

	sql := `
	INSERT INTO articles (author_id, title, slug, description, body, favorites_count) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at
	`
	args := []any{
		article.Author.Id,
		article.Title,
		article.Slug,
		article.Description,
		article.Body,
		article.FavoritesCount,
	}

	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(&article.Id, &article.CreatedAt, &article.UpdatedAt)
	return errx.ErrxOf(err)
}

// ArticleGetBySlugDaf implements a stereotype instance of type
// fs.ArticleGetBySlugDafT.
var ArticleGetBySlugDaf fs.ArticleGetBySlugDafT = func(
	ctx context.Context,
	slug string,
) (model.Article, error) {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return model.Article{}, errx.ErrxOf(err)
	}

	sql := `
	SELECT * FROM articles WHERE slug = $1
	`
	rows, err := tx.Query(ctx, sql, slug)
	if err != nil {
		return model.Article{}, errx.ErrxOf(err)
	}
	var article model.Article
	err = pgxscan.ScanOne(&article, rows)
	if err != nil {
		return model.Article{}, errx.ErrxOf(err)
	}
	return article, nil
}

// ArticleUpdateDaf implements a stereotype instance of type
// fs.ArticleUpdateDafT.
var ArticleUpdateDaf fs.ArticleUpdateDafT = func(
	ctx context.Context,
	article *model.Article,
) error {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return errx.ErrxOf(err)
	}
	sql := `
	UPDATE articles 
	SET title = $1, description = $2, body = $3, updated_at = NOW() 
	WHERE id = $4 AND updated_at = $5
	RETURNING updated_at
	`
	args := []interface{}{
		article.Title,
		article.Description,
		article.Body,
		article.Id,
		article.UpdatedAt,
	}
	row := tx.QueryRow(ctx, sql, args...)
	err = row.Scan(&article.UpdatedAt)
	return errx.ErrxOf(err)
}

// ArticleDeleteDaf implements a stereotype instance of type
// fs.ArticleDeleteDafT.
var ArticleDeleteDaf fs.ArticleDeleteDafT = func(
	ctx context.Context,
	slug string,
) error {
	tx, err := dbpgx.GetCtxTx(ctx)
	if err != nil {
		return errx.ErrxOf(err)
	}
	sql := `
	DELETE FROM articles
	WHERE slug = $1
	`
	_, err = tx.Exec(ctx, sql, slug)
	return errx.ErrxOf(err)
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
