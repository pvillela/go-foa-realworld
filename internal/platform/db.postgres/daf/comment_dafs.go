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

// CommentCreateDafI implements a stereotype instance of type
// CommentCreateDafT.
var CommentCreateDafI CommentCreateDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	comment *model.Comment,
) error {
	sql := `
	INSERT INTO comments (article_id, author_id, body) 
	VALUES ($1, $2, $3) 
	RETURNING id, created_at
	`
	args := []any{
		comment.ArticleId,
		comment.AuthorId,
		comment.Body,
	}

	row := tx.QueryRow(ctx, sql, args...)
	err := row.Scan(&comment.Id, &comment.CreatedAt)
	return errx.ErrxOf(err)
}

// CommentsGetBySlugDafI implements a stereotype instance of type
// CommentGetBySlugDafT.
var CommentsGetBySlugDafI CommentsGetBySlugDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	slug string,
) ([]model.Comment, error) {
	sql := `
	SELECT c.* FROM comments c
	LEFT JOIN articles a ON c.author_id = a.id
	WHERE a.slug = $1
	`
	args := []any{slug}

	comments, err := dbpgx.ReadMany[model.Comment](ctx, tx, sql, -1, -1, args...)
	return comments, err
}

// CommentDeleteDafI implements a stereotype instance of type
// CommentDeleteDafT.
var CommentDeleteDafI CommentDeleteDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	id uint,
) error {
	sql := `
	DELETE FROM comments
	WHERE id = $1
	`
	_, err := tx.Exec(ctx, sql, id)
	return errx.ErrxOf(err)
}
