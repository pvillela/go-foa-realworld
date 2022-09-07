/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/pvillela/go-foa-realworld/experimental/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/experimental/arch/errx"
	"github.com/pvillela/go-foa-realworld/experimental/arch/util"
	"github.com/pvillela/go-foa-realworld/experimental/daf"
	"github.com/pvillela/go-foa-realworld/experimental/model"
	"github.com/pvillela/go-foa-realworld/experimental/rpc"
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
	errx.PanicOnError(err)

	for i, _ := range tags {
		err := daf.TagCreateDaf(ctx, tx, &tags[i])
		errx.PanicOnError(err)
		fmt.Println("tag from Create:", tags[i], "\n")

		err = daf.TagAddToArticleDaf(ctx, tx, tags[i], articles[1])
		errx.PanicOnError(err)
	}

	tagsFromDb, err := daf.TagsGetAllDaf(ctx, tx)
	fmt.Println("Tags from database:", tagsFromDb, "\n")

	// Try to insert same tag again
	{
		err = tx.Commit(ctx)
		errx.PanicOnError(err)
		tx, err = db.BeginTx(ctx)
		errx.PanicOnError(err)

		err = daf.TagCreateDaf(ctx, tx, &tags[0])
		fmt.Println("Duplicate tag insert:", err)
		fmt.Printf("pgconn.PgError: %+v", *errors.Unwrap(err).(*pgconn.PgError))
		fmt.Println("SqlState(err):", dbpgx.SqlState(err), "\n")

		err = tx.Commit(ctx)
		log.Debug(err, "\n\n")
		tx, err = db.BeginTx(ctx)
		errx.PanicOnError(err)
	}

	currUserId := users[0].Id

	{
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
	}

	{
		// Add some more tags to database.
		err = daf.TagsAddNewDaf(ctx, tx, []string{"ZZZ", "FOOTAG", "WWW"})
		errx.PanicOnError(err)

		tagsFromDb, err = daf.TagsGetAllDaf(ctx, tx)
		fmt.Println("Tags from database:", tagsFromDb, "\n")

		// Add the new tags to an article.
		err = daf.TagsAddToArticleDaf(ctx, tx, []string{"ZZZ", "FOOTAG", "WWW"}, articles[1])
		errx.PanicOnError(err)
	}

	{
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
	}

	{
		criteria := rpc.ArticleCriteria{
			Tag:         util.PointerFromValue("FOOTAG"),
			Author:      nil,
			FavoritedBy: nil,
			Limit:       nil,
			Offset:      nil,
		}
		articlePluses, err := daf.ArticlesListDaf(ctx, tx, currUserId, criteria)
		errx.PanicOnError(err)
		fmt.Println("\narticlesListDaf - by tag:", articlePluses, "\n")
	}

	{
		articleFromDb, err := daf.ArticleGetBySlugDaf(ctx, tx, currUserId, articles[1].Slug)
		errx.PanicOnError(err)
		fmt.Println("\nArticleGetBySlugDaf:", articleFromDb, "\n")
	}

	err = tx.Commit(ctx)
	errx.PanicOnError(err)
}
