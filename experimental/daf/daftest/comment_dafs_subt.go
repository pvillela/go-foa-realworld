/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/experimental/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/experimental/arch/db/dbpgx/dbpgxtest"
	"github.com/pvillela/go-foa-realworld/experimental/arch/errx"
	"github.com/pvillela/go-foa-realworld/experimental/arch/util"
	"github.com/pvillela/go-foa-realworld/experimental/model"
	"github.com/pvillela/go-foa-realworld/experimental/daf"
	"github.com/stretchr/testify/assert"
	"testing"
)

var commentDafsSubt = dbpgxtest.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {

	// Create comments.

	commentsDeleteAll(ctx, tx) // cleanup to make test idempotent because comments are cumulative
	mdb.CommentDeleteAll()

	type commentSourceT struct {
		username string
		slug     string
		body     string
	}

	commentSources := []commentSourceT{
		{
			username: username1,
			slug:     slug1,
			body:     "I thought it was a great article.",
		},
		{
			username: username3,
			slug:     slug1,
			body:     "What a terrible article.",
		},
	}

	makeComment := func(src commentSourceT) (string, string, model.Comment) {
		comment := model.Comment{
			ArticleId: mdb.ArticleGetBySlug(src.slug).Id,
			AuthorId:  mdb.UserGetByName(src.username).Id,
			Body:      util.PointerFromValue(src.body),
		}
		return src.username, src.slug, comment
	}

	for _, csrc := range commentSources {
		username, slug, comment := makeComment(csrc)
		err := daf.CommentCreateDaf(ctx, tx, &comment)
		assert.NoError(t, err)
		mdb.CommentInsert(username, slug, comment)
	}

	// Tests

	{
		msg := "Get comments for article with comments."

		slug := slug1

		returned, err := daf.CommentsGetBySlugDaf(ctx, tx, slug)
		assert.NoError(t, err)

		expected := mdb.CommentGetAllBySlug(slug)

		assert.ElementsMatch(t, expected, returned, msg)
	}

	{
		msg := "Get comments for article without comments."

		slug := slug2

		returned, err := daf.CommentsGetBySlugDaf(ctx, tx, slug)
		assert.NoError(t, err)

		expected := mdb.CommentGetAllBySlug(slug)

		assert.Equal(t, len(returned), 0,
			"this test requires the returned slice to be empty")
		assert.ElementsMatch(t, expected, returned, msg)
	}

	{
		msg := "Deletion of comment by author."

		currUsername := username1
		commentAuthorUsername := currUsername
		slug := slug1

		currUser := mdb.UserGetByName(currUsername)
		comment := mdb.CommentGetAllForKey(commentAuthorUsername, slug)[0]

		err := daf.CommentDeleteDaf(ctx, tx, comment.Id, comment.ArticleId, currUser.Id)
		assert.NoError(t, err)

		// Sync in-memory data
		mdb.CommentDelete(currUsername, slug, comment.Id)

		{
			returned := commentsGetAll(ctx, tx)
			expected := mdb.CommentGetAll()

			assert.ElementsMatch(t, expected, returned, msg)
		}
	}

	{
		msg := "Deletion of comment by non-author. This does not invalidate the transaction."

		currUsername := username1
		commentAuthorUsername := username3
		slug := slug1

		currUser := mdb.UserGetByName(currUsername)
		comment := mdb.CommentGetAllForKey(commentAuthorUsername, slug)[0]

		err := daf.CommentDeleteDaf(ctx, tx, comment.Id, comment.ArticleId, currUser.Id)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg)

		{
			returned := commentsGetAll(ctx, tx)
			expected := mdb.CommentGetAll()

			assert.ElementsMatch(t, expected, returned, msg)
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

func commentsDeleteAll(ctx context.Context, tx pgx.Tx) {
	sql := `
	DELETE FROM comments
	`
	_, err := tx.Exec(ctx, sql)
	errx.PanicOnError(err)
}
