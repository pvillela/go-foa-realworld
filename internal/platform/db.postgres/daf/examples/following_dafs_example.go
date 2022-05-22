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

func followingDafsExample(ctx context.Context, db dbpgx.Db) {
	fmt.Println("********** followingDafsExample **********\n")

	tx, err := db.BeginTx(ctx)
	util.PanicOnError(err)

	followings := []model.Following{
		{
			FollowerID: users[0].Id,
			FolloweeID: users[1].Id,
		},
		{
			FollowerID: users[2].Id,
			FolloweeID: users[1].Id,
		},
	}

	for i, _ := range followings {
		err := daf.FollowingCreateDafI(ctx, tx, followings[i].FollowerID, followings[i].FolloweeID)
		util.PanicOnError(err)
	}

	{
		currUserId := users[0].Id

		articlePluses, err := daf.ArticlesFeedDafI(ctx, tx, currUserId, nil, nil)
		util.PanicOnError(err)
		fmt.Println("\nArticlesFeedDaf:", articlePluses, "\n")
	}

	{
		currUserId := users[2].Id

		articlePluses, err := daf.ArticlesFeedDafI(ctx, tx, currUserId, nil, nil)
		util.PanicOnError(err)
		fmt.Println("\nArticlesFeedDaf:", articlePluses, "\n")
	}

	{
		_, err := daf.FollowingDeleteDafI(ctx, tx, followings[1].FollowerID, followings[1].FolloweeID)
		util.PanicOnError(err)

		currUserId := users[2].Id

		articlePluses, err := daf.ArticlesFeedDafI(ctx, tx, currUserId, nil, nil)
		util.PanicOnError(err)
		fmt.Println("\nArticlesFeedDaf:", articlePluses, "\n")
	}

	err = tx.Commit(ctx)
	util.PanicOnError(err)
}
