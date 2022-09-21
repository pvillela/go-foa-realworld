/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

const (
	slug1 = "anintsubj"
	slug2 = "adullsubj"
)

var articles = []model.Article{
	{
		Title:       "An interesting subject",
		Slug:        slug1,
		Description: "Story about an interesting subject.",
		Body:        util.PointerOf("I met this interesting subject a long time ago."),
	},
	{
		Title:       "A dull story",
		Slug:        slug2,
		Description: "Narrative about something dull.",
		Body:        util.PointerOf("This is so dull, bla, bla, bla."),
	},
}

func articleDafsExample(ctx context.Context, db dbpgx.Db) {
	fmt.Println("********** articleDafsExample **********\n")

	tx, err := db.BeginTx(ctx)
	errx.PanicOnError(err)

	for i, _ := range articles {
		articles[i].AuthorId = users[1].Id
		err := daf.ArticleCreateDaf(ctx, tx, &articles[i])
		errx.PanicOnError(err)
		fmt.Println("article from Create:", articles[i], "\n")
	}

	err = tx.Commit(ctx)
	errx.PanicOnError(err)

	tx, err = db.BeginTx(ctx)
	errx.PanicOnError(err)

	currUserId := users[0].Id

	criteria := rpc.ArticleCriteria{
		Tag:         nil,
		Author:      &users[1].Username,
		FavoritedBy: nil,
		Limit:       nil,
		Offset:      nil,
	}
	articlePluses, err := daf.ArticlesListDaf(ctx, tx, currUserId, criteria)
	errx.PanicOnError(err)
	fmt.Println("\narticlesListDaf - by author:", articlePluses, "\n")

	criteria = rpc.ArticleCriteria{
		Tag:         util.PointerOf("FOOTAG"),
		Author:      nil,
		FavoritedBy: nil,
		Limit:       nil,
		Offset:      nil,
	}
	articlePluses, err = daf.ArticlesListDaf(ctx, tx, currUserId, criteria)
	errx.PanicOnError(err)
	fmt.Println("\narticlesListDaf - by tag:", articlePluses, "\n")

	articleFromDb, err := daf.ArticleGetBySlugDaf(ctx, tx, currUserId, articles[1].Slug)
	errx.PanicOnError(err)
	fmt.Println("\nArticleGetBySlugDaf:", articleFromDb, "\n")

	pArticle := &articles[0]
	pArticle.Title = "A very interesting subject"
	err = daf.ArticleUpdateDaf(ctx, tx, pArticle)
	errx.PanicOnError(err)
	fmt.Println("ArticleUpdateDaf:", pArticle, "\n")

	articlePluses, err = daf.ArticlesFeedDaf(ctx, tx, currUserId, nil, nil)
	errx.PanicOnError(err)
	fmt.Println("\nArticlesFeedDaf:", articlePluses, "\n")

	err = tx.Commit(ctx)
	errx.PanicOnError(err)
}
