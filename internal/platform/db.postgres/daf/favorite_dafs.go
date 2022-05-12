/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
)

// FavoriteCreateDaf is the instance of the DAF stereotype that
// associates an article with a user that favors it.
func FavoriteCreateDaf(ctx context.Context, tx pgx.Tx, articleId uint, userId uint) error {
	sql := `
	INSERT INTO favorites (article_id, user_id)
	VALUES ($1, $2)
	`
	args := []any{articleId, userId}
	_, err := tx.Exec(ctx, sql, args...)
	return errx.ErrxOf(err)
}

// FavoriteDeleteDaf is the instance of the DAF stereotype that
// disaassociates an article from a user that favors it.
func FavoriteDeleteDaf(ctx context.Context, tx pgx.Tx, articleId uint, userId uint) error {
	sql := `
	DELETE FROM favorites
	WHERE article_id = $1 AND user_id = $2
	`
	args := []any{articleId, userId}
	_, err := tx.Exec(ctx, sql, args...)
	return errx.ErrxOf(err)
}
