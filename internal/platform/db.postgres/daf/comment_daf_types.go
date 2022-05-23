/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// CommentsGetBySlugDafT is the type of the stereotype instance for the DAF that
// retrieves all comments for the article with a given slug.
type CommentsGetBySlugDafT = func(ctx context.Context, tx pgx.Tx, slug string) ([]model.Comment, error)

// CommentCreateDafT is the type of the stereotype instance for the DAF that
// creates a comment for an article.
type CommentCreateDafT = func(ctx context.Context, tx pgx.Tx, comment *model.Comment) error

// CommentDeleteDafT is the type of the stereotype instance for the DAF that
// deletes a comment.
type CommentDeleteDafT = func(ctx context.Context, tx pgx.Tx, id uint) error
