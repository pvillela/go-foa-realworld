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
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// CommentCreateDaf implements a stereotype instance of type
// CommentCreateDafT.
var CommentCreateDaf CommentCreateDafT = func(
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
	if kind := dbpgx.ClassifyError(err); kind != nil {
		return kind.Make(err, "")
	}
	return nil
}

// CommentsGetBySlugDaf implements a stereotype instance of type
// CommentGetBySlugDafT.
var CommentsGetBySlugDaf CommentsGetBySlugDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	slug string,
) ([]model.Comment, error) {
	sql := `
	SELECT c.* FROM comments c
	LEFT JOIN articles a ON c.article_id = a.id
	WHERE a.slug = $1
	`
	args := []any{slug}

	comments, err := dbpgx.ReadMany[model.Comment](ctx, tx, sql, -1, -1, args...)
	return comments, err
}

// CommentDeleteDaf implements a stereotype instance of type
// CommentDeleteDafT.
var CommentDeleteDaf CommentDeleteDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	commentId uint,
	articleId uint,
	authorId uint,
) error {
	sql := `
	DELETE FROM comments
	WHERE id = $1 AND article_id = $2 AND author_id = $3
	`
	args := []any{
		commentId,
		articleId,
		authorId,
	}

	c, err := tx.Exec(ctx, sql, args...)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		return kind.Make(err, "")
	}
	if c.RowsAffected() == 0 {
		return dbpgx.DbErrRecordNotFound.Make(nil, bf.ErrMsgCommentNotFound)
	}

	return nil
}
