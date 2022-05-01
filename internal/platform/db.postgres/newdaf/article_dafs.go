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
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

const (
	limitDefault  = 20
	offsetDefault = 0
)

// ArticleCreateDaf implements a stereotype instance of type
// ArticleCreateDafT.
var ArticleCreateDaf ArticleCreateDafT = func(
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
// ArticleGetBySlugDafT.
var ArticleGetBySlugDaf ArticleGetBySlugDafT = func(
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
	defer rows.Close()

	var article model.Article
	err = pgxscan.ScanOne(&article, rows)
	if err != nil {
		return model.Article{}, errx.ErrxOf(err)
	}
	return article, nil
}

// ArticleUpdateDaf implements a stereotype instance of type
// ArticleUpdateDafT.
var ArticleUpdateDaf ArticleUpdateDafT = func(
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
// ArticleDeleteDafT.
var ArticleDeleteDaf ArticleDeleteDafT = func(
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

// ArticlesFeedDaf implements a stereotype instance of type
// ArticlesFeedDafT.
var ArticlesFeedDaf ArticlesFeedDafT = func(
	ctx context.Context,
	in rpc.ArticlesFeedIn,
) (rpc.ArticlesOut, error) {
	sql := `
	
	`
}

// ArticlesListDafT implements a stereotype instance of type
// ArticlesListDafTT.
var ArticlesListDaf ArticlesListDafT = func(
	ctx context.Context,
	in rpc.ArticlesListIn,
) (rpc.ArticlesOut, error) {
	
}
