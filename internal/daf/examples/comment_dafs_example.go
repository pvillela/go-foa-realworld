/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

func commentDafsExample(ctx context.Context, db dbpgx.Db) {
	fmt.Println("********** CommentDafsExample **********")

	tx, err := db.BeginTx(ctx)
	errx.PanicOnError(err)

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
		err := daf.CommentCreateDaf(ctx, tx, &comments[i])
		errx.PanicOnError(err)
	}

	{
		commentsFromDb, err := daf.CommentsGetBySlugDaf(ctx, tx, articles[0].Slug)
		errx.PanicOnError(err)
		_, _ = spew.Printf("\nCommentsGetBySlugDaf: %v\n", commentsFromDb)
	}

	{
		// Deletion of article by author.
		currUser := users[0]
		comment := comments[0]

		fmt.Printf(
			"\narticle deletion - currUser ID: %v, comment author ID: %v:\n",
			currUser.Id,
			comment.AuthorId,
		)
		err := daf.CommentDeleteDaf(ctx, tx, comment.Id, comment.ArticleId, currUser.Id)
		errx.PanicOnError(err)

		/////////////
		// IMPORTANT: removed comments[0] from memory copy to sync with database
		comments = comments[1:]

		{
			commentsFromDb, err := daf.CommentsGetBySlugDaf(ctx, tx, articles[0].Slug)
			errx.PanicOnError(err)
			_, _ = spew.Printf("\nCommentsGetBySlugDaf: %v\n", commentsFromDb)

			commentsFromDb = commentsGetAll(ctx, tx)
			_, _ = spew.Printf("\ncommentsGetAll: %v\n", commentsFromDb)
		}
	}

	{
		// Deletion of article by non-author. This does not invalidate the transaction.
		currUser := users[0]
		comment := comments[0] // is the old comments[1] due to previous delete operation

		fmt.Printf(
			"\narticle deletion - currUser ID: %v, comment author ID: %v:\n",
			currUser.Id,
			comment.AuthorId,
		)
		err := daf.CommentDeleteDaf(ctx, tx, comment.Id, comment.ArticleId, currUser.Id)
		fmt.Println("err:", err)

		{
			commentsFromDb, err := daf.CommentsGetBySlugDaf(ctx, tx, articles[0].Slug)
			errx.PanicOnError(err)
			_, _ = spew.Printf("\nCommentsGetBySlugDaf: %v\n", commentsFromDb)

			commentsFromDb = commentsGetAll(ctx, tx)
			_, _ = spew.Printf("\ncommentsGetAll: %v\n", commentsFromDb)
		}
	}

	err = tx.Commit(ctx)
	errx.PanicOnError(err)
}

func commentsGetAll(ctx context.Context, tx pgx.Tx) []model.Comment {
	sql := `
	SELECT * FROM comments
	`
	comments, err := dbpgx.ReadMany[model.Comment](ctx, tx, sql, -1, -1)
	errx.PanicOnError(err)

	return comments
}
