/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// TagGetAllDaf implements a stereotype instance of type
// TagGetAllDafT.
var TagGetAllDaf TagGetAllDafT = func(ctx context.Context, tx pgx.Tx) ([]string, error) {
	mainSql := `
	SELECT * FROM tags
	`
	return dbpgx.ReadMany[string](ctx, tx, mainSql, -1, -1)
}

// TagCreateDaf implements a stereotype instance of type
// TagCreateDafT.
var TagCreateDaf TagCreateDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	tag *model.Tag,
) error {
	sql := `
	INSERT INTO tags (name)
	VALUES ($1)
	RETURNING id
	`
	args := []any{tag.Name}
	row := tx.QueryRow(ctx, sql, args...)
	err := row.Scan(&tag.Id)
	return errx.ErrxOf(err)
}

// TagAddToArticle implements a stereotype instance of type
// TagAddToArticleT.
var TagAddToArticle TagAddToArticleT = func(
	ctx context.Context,
	tx pgx.Tx,
	tag model.Tag,
	article model.Article,
) error {
	sql := `
	INSERT INTO article_tags (article_id, tag_id)
	VALUES ($1, $2)
	`
	args := []any{article.Id, tag.Id}
	_, err := tx.Exec(ctx, sql, args...)
	return errx.ErrxOf(err)
}