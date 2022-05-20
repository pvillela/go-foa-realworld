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
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

var articles = []model.Article{
	{
		Title:       "An interesting subject",
		Slug:        "anintsubj",
		Description: "Story about an interesting subject.",
		Body:        util.PointerFromValue("I met this interesting subject a long time ago."),
	},
	{
		Title:       "A dull story",
		Slug:        "adullsubj",
		Description: "Narrative about something dull.",
		Body:        util.PointerFromValue("This is so dull, bla, bla, bla."),
	},
}

func articleDafsExample(ctx context.Context, db dbpgx.Db) {
	fmt.Println("********** articleDafsExample **********\n")

	tx, err := db.BeginTx(ctx)
	util.PanicOnError(err)

	for i, _ := range articles {
		articles[i].AuthorId = users[1].Id
		err := daf.ArticleCreateDafI(ctx, tx, &articles[i])
		util.PanicOnError(err)
		fmt.Println("article from Create:", articles[i], "\n")
	}

	err = tx.Commit(ctx)
	util.PanicOnError(err)

	tx, err = db.BeginTx(ctx)
	util.PanicOnError(err)

	currUserId := users[0].Id

	criteria := model.ArticleCriteria{
		Tag:         nil,
		Author:      &users[1].Username,
		FavoritedBy: nil,
		Limit:       nil,
		Offset:      nil,
	}
	articlePluses, err := daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
	util.PanicOnError(err)
	fmt.Println("\narticlesListDaf - by author:", articlePluses, "\n")

	criteria = model.ArticleCriteria{
		Tag:         util.PointerFromValue("FOOTAG"),
		Author:      nil,
		FavoritedBy: nil,
		Limit:       nil,
		Offset:      nil,
	}
	articlePluses, err = daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
	util.PanicOnError(err)
	fmt.Println("\narticlesListDaf - by tag:", articlePluses, "\n")

	articleFromDb, err := daf.ArticleGetBySlugDafI(ctx, tx, currUserId, articles[1].Slug)
	util.PanicOnError(err)
	fmt.Println("\nArticleGetBySlugDaf:", articleFromDb, "\n")

	pArticle := &articles[0]
	pArticle.Title = "A very interesting subject"
	err = daf.ArticleUpdateDafI(ctx, tx, pArticle)
	util.PanicOnError(err)
	fmt.Println("ArticleUpdateDaf:", pArticle, "\n")

	articlePluses, err = daf.ArticlesFeedDafI(ctx, tx, currUserId, nil, nil)
	util.PanicOnError(err)
	fmt.Println("\nArticlesFeedDaf:", articlePluses, "\n")

	err = tx.Commit(ctx)
	util.PanicOnError(err)
}
