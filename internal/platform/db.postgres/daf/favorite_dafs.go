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
	"github.com/pvillela/go-foa-realworld/internal/bf"
)

// FavoriteCreateDafI is the instance of the DAF stereotype that
// associates an article with a user that favors it.
var FavoriteCreateDafI FavoriteCreateDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	articleId uint,
	userId uint,
) error {
	sql := `
	INSERT INTO favorites (article_id, user_id)
	VALUES ($1, $2)
	`
	args := []any{articleId, userId}
	_, err := tx.Exec(ctx, sql, args...)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind == dbpgx.DbErrUniqueViolation {
			return kind.Make(err, bf.ErrMsgArticleAlreadyFavorited, articleId)
		}
		return kind.Make(err, "")
	}
	return nil
}

// FavoriteDeleteDafI is the instance of the DAF stereotype that
// disaassociates an article from a user that favors it.
var FavoriteDeleteDafI FavoriteDeleteDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	articleId uint,
	userId uint,
) error {
	sql := `
	DELETE FROM favorites
	WHERE article_id = $1 AND user_id = $2
	`
	args := []any{articleId, userId}
	c, err := tx.Exec(ctx, sql, args...)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		return kind.Make(err, "")
	}
	if c.RowsAffected() == 0 {
		return dbpgx.DbErrRecordNotFound.Make(nil, bf.ErrMsgArticleWasNotFavorited, articleId)
	}

	return nil
}
