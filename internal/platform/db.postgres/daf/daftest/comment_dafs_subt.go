/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/stretchr/testify/assert"
	"testing"
)

var commentDafsSubt = dbpgx.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {

	// Create comments.

	comments := []model.Comment{
		{
			ArticleId: articles[0].Id,
			AuthorId:  users[0].Id,
			Body:      util.PointerFromValue("I thought it was a great article."),
		},
		{
			ArticleId: articles[0].Id,
			AuthorId:  users[2].Id,
			Body:      util.PointerFromValue("What a terrible article."),
		},
	}

	for i, _ := range comments {
		err := daf.CommentCreateDafI(ctx, tx, &comments[i])
		errx.PanicOnError(err)
	}

	// Tests

	{
		msg := "Get comments for article with comments."

		returned, err := daf.CommentsGetBySlugDafI(ctx, tx, articles[0].Slug)
		errx.PanicOnError(err)

		assert.Equal(t, comments, returned, msg)
	}

	{
		msg := "Get comments for article without comments."

		returned, err := daf.CommentsGetBySlugDafI(ctx, tx, articles[1].Slug)
		errx.PanicOnError(err)

		var expected []model.Comment

		assert.Equal(t, expected, returned, msg)
	}

	{
		msg := "Deletion of article by author."
		currUser := users[0]
		comment := comments[0]

		err := daf.CommentDeleteDafI(ctx, tx, comment.Id, comment.ArticleId, currUser.Id)
		errx.PanicOnError(err)

		/////////////
		// IMPORTANT: removed comments[0] from memory copy to sync with database
		comments = comments[1:]

		{
			returned := commentsGetAll(ctx, tx)

			assert.ElementsMatch(t, comments, returned, msg)
		}
	}

	{
		msg := "Deletion of article by non-author. This does not invalidate the transaction."
		currUser := users[0]
		comment := comments[0] // is the old comments[1] due to previous delete operation

		err := daf.CommentDeleteDafI(ctx, tx, comment.Id, comment.ArticleId, currUser.Id)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg)

		{
			returned := commentsGetAll(ctx, tx)

			assert.ElementsMatch(t, comments, returned, msg)
		}
	}
})

func commentsGetAll(ctx context.Context, tx pgx.Tx) []model.Comment {
	sql := `
	SELECT * FROM comments
	`
	comments, err := dbpgx.ReadMany[model.Comment](ctx, tx, sql, -1, -1)
	errx.PanicOnError(err)

	return comments
}
