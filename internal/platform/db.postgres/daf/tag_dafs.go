/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	log "github.com/sirupsen/logrus"
	"strings"
)

// TagsGetAllDafI implements a stereotype instance of type
// TagsGetAllDafT.
var TagsGetAllDafI TagsGetAllDafT = func(ctx context.Context, tx pgx.Tx) ([]model.Tag, error) {
	mainSql := `
	SELECT * FROM tags
	`
	return dbpgx.ReadMany[model.Tag](ctx, tx, mainSql, -1, -1)
}

// TagCreateDafI implements a stereotype instance of type
// TagCreateDafT.
var TagCreateDafI TagCreateDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	tag *model.Tag,
) error {
	sql := `
	INSERT INTO tags (name)
	VALUES ($1)
	RETURNING id
	`
	args := []any{strings.ToUpper(tag.Name)}
	row := tx.QueryRow(ctx, sql, args...)
	err := row.Scan(&tag.Id)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind == dbpgx.DbErrUniqueViolation {
			return kind.Make(err, bf.ErrMsgTagNameAlreadyExists, tag.Name)
		}
		return kind.Make(err, "")
	}
	return nil
}

// TagsAddNewDafI implements a stereotype instance of type
// TagsAddNewDafT.
var TagsAddNewDafI TagsAddNewDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	names []string,
) error {
	if len(names) == 0 {
		return nil
	}

	preSql := `
	INSERT INTO tags (name)
	SELECT x.name
	FROM (
		VALUES
			%v
	) x (name)
	WHERE NOT EXISTS (
		SELECT 1
		FROM tags t
		WHERE t.name = x.name
	)
	`
	var values []string
	for _, name := range names {
		values = append(values, fmt.Sprintf("('%v')", name))
	}
	valueString := strings.Join(values, ", ")
	sql := fmt.Sprintf(preSql, valueString)
	log.Debug("sql:", sql)

	_, err := tx.Exec(ctx, sql)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		return kind.Make(err, "")
	}

	return nil
}

// TagAddToArticleDafI implements a stereotype instance of type
// TagAddToArticleDafT.
var TagAddToArticleDafI TagAddToArticleDafT = func(
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
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind == dbpgx.DbErrUniqueViolation {
			return kind.Make(err, bf.ErrMsgTagOnArticleAlreadyExists, tag.Name, article.Slug)
		}
		return kind.Make(err, "")
	}
	return nil
}

// TagsAddToArticleDafI implements a stereotype instance of type
// TagsAddToArticleDafT.
var TagsAddToArticleDafI TagsAddToArticleDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	names []string,
	article model.Article,
) error {
	preSql := `
	INSERT INTO article_tags (article_id, tag_id)
	SELECT $1, t.id
	FROM tags t
	WHERE t.name IN (%v)
	AND NOT EXISTS (
			SELECT 1
			FROM article_tags at
			WHERE at.tag_id = t.id
	)
	`
	var values []string
	for _, name := range names {
		values = append(values, fmt.Sprintf("'%v'", name))
	}
	valueString := strings.Join(values, ", ")
	sql := fmt.Sprintf(preSql, valueString)
	log.Debug("sql:", sql)

	_, err := tx.Exec(ctx, sql, article.Id)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		return kind.Make(err, "")
	}

	return nil
}
