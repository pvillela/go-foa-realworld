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
		err := daf.CommentCreateDafI(ctx, tx, &comments[i])
		errx.PanicOnError(err)
	}

	{
		comments, err := daf.CommentsGetBySlugDafI(ctx, tx, articles[0].Slug)
		errx.PanicOnError(err)
		_, _ = spew.Printf("\nCommentsGetBySlugDafI: %v\n", comments)
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
		err := daf.CommentDeleteDafI(ctx, tx, comment.Id, comment.ArticleId, currUser.Id)
		errx.PanicOnError(err)
	}

	{
		comments, err := daf.CommentsGetBySlugDafI(ctx, tx, articles[0].Slug)
		errx.PanicOnError(err)
		_, _ = spew.Printf("\nCommentsGetBySlugDafI: %v\n", comments)
	}

	{
		// Deletion of article by non-author. This does not invalidate the transaction.
		currUser := users[2]
		comment := comments[0]

		fmt.Printf(
			"\narticle deletion - currUser ID: %v, comment author ID: %v:\n",
			currUser.Id,
			comment.AuthorId,
		)
		err := daf.CommentDeleteDafI(ctx, tx, comment.Id, comment.ArticleId, currUser.Id)
		fmt.Println("err:", err)
	}

	{
		comments, err := daf.CommentsGetBySlugDafI(ctx, tx, articles[0].Slug)
		errx.PanicOnError(err)
		_, _ = spew.Printf("\nCommentsGetBySlugDafI: %v\n", comments)
	}

	err = tx.Commit(ctx)
	errx.PanicOnError(err)
}
