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
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

func followingDafsExample(ctx context.Context, db dbpgx.Db) {
	fmt.Println("********** followingDafsExample **********\n")

	tx, err := db.BeginTx(ctx)
	errx.PanicOnError(err)

	followings := []model.Following{
		{
			FollowerId: users[0].Id,
			FolloweeId: users[1].Id,
		},
		{
			FollowerId: users[2].Id,
			FolloweeId: users[1].Id,
		},
	}

	for i, _ := range followings {
		_, err := daf.FollowingCreateDaf(ctx, tx, followings[i].FollowerId, followings[i].FolloweeId)
		errx.PanicOnError(err)
	}

	{
		currUserId := users[0].Id

		articlePluses, err := daf.ArticlesFeedDaf(ctx, tx, currUserId, nil, nil)
		errx.PanicOnError(err)
		fmt.Println("\nArticlesFeedDaf:", articlePluses, "\n")
	}

	{
		currUserId := users[2].Id

		articlePluses, err := daf.ArticlesFeedDaf(ctx, tx, currUserId, nil, nil)
		errx.PanicOnError(err)
		fmt.Println("\nArticlesFeedDaf:", articlePluses, "\n")
	}

	{
		err := daf.FollowingDeleteDaf(ctx, tx, followings[1].FollowerId, followings[1].FolloweeId)
		errx.PanicOnError(err)

		currUserId := users[2].Id

		articlePluses, err := daf.ArticlesFeedDaf(ctx, tx, currUserId, nil, nil)
		errx.PanicOnError(err)
		fmt.Println("\nArticlesFeedDaf:", articlePluses, "\n")
	}

	err = tx.Commit(ctx)
	errx.PanicOnError(err)
}
