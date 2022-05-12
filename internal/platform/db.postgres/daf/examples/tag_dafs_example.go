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
	log "github.com/sirupsen/logrus"
)

var tags = []model.Tag{
	{
		Name: "FOOTAG",
	},
	{
		Name: "BARTAG",
	},
}

func tagDafsExample(ctx context.Context, db dbpgx.Db) {
	fmt.Println("********** tagDafsExample **********\n")

	tx, err := db.BeginTx(ctx)
	util.PanicOnError(err)

	for i, _ := range tags {
		err := daf.TagCreateDaf(ctx, tx, &tags[i])
		util.PanicOnError(err)
		fmt.Println("tag from Create:", tags[i], "\n")

		err = daf.TagAddToArticle(ctx, tx, tags[i], articles[1])
		util.PanicOnError(err)
	}

	// Try to insert same tag again
	{
		err = tx.Commit(ctx)
		util.PanicOnError(err)
		tx, err = db.BeginTx(ctx)
		util.PanicOnError(err)

		err = daf.TagCreateDaf(ctx, tx, &tags[0])
		fmt.Println("Duplicate tag insert:", err)
		fmt.Println("SqlState(err):", dbpgx.SqlState(err), "\n")

		err = tx.Commit(ctx)
		log.Debug(err, "\n\n")
		tx, err = db.BeginTx(ctx)
		util.PanicOnError(err)
	}

	currUserId := users[0].Id

	criteria := model.ArticleCriteria{
		Tag:         nil,
		Author:      &users[1].Username,
		FavoritedBy: nil,
		Limit:       nil,
		Offset:      nil,
	}
	articlePluses, err := daf.ArticlesListDaf(ctx, tx, currUserId, criteria)
	util.PanicOnError(err)
	fmt.Println("\narticlesListDaf - by author:", articlePluses, "\n")

	criteria = model.ArticleCriteria{
		Tag:         util.PointerFromValue("FOOTAG"),
		Author:      nil,
		FavoritedBy: nil,
		Limit:       nil,
		Offset:      nil,
	}
	articlePluses, err = daf.ArticlesListDaf(ctx, tx, currUserId, criteria)
	util.PanicOnError(err)
	fmt.Println("\narticlesListDaf - by tag:", articlePluses, "\n")

	articleFromDb, err := daf.ArticleGetBySlugDaf(ctx, tx, currUserId, articles[1].Slug)
	util.PanicOnError(err)
	fmt.Println("\nArticleGetBySlugDaf:", articleFromDb, "\n")

	err = tx.Commit(ctx)
	util.PanicOnError(err)
}
