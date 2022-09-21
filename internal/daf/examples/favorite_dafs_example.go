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

func favoriteDafsExample(ctx context.Context, db dbpgx.Db) {
	fmt.Println("********** FavoriteDafsExample **********\n")

	tx, err := db.BeginTx(ctx)
	errx.PanicOnError(err)

	Favorites := []model.Favorite{
		{
			ArticleId: articles[0].Id,
			UserId:    users[0].Id,
		},
		{
			ArticleId: articles[0].Id,
			UserId:    users[2].Id,
		},
	}

	for i, _ := range Favorites {
		err := daf.FavoriteCreateDaf(ctx, tx, Favorites[i].ArticleId, Favorites[i].UserId)
		errx.PanicOnError(err)
	}

	{
		currUserId := users[0].Id

		criteria := rpc.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: util.PointerOf(users[2].Username),
			Limit:       nil,
			Offset:      nil,
		}
		articlePluses, err := daf.ArticlesListDaf(ctx, tx, currUserId, criteria)
		errx.PanicOnError(err)
		fmt.Println("\narticlesListDaf - favoritedBy:", articlePluses, "\n")
	}

	{
		currUserId := users[2].Id

		criteria := rpc.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: util.PointerOf(users[0].Username),
			Limit:       nil,
			Offset:      nil,
		}
		articlePluses, err := daf.ArticlesListDaf(ctx, tx, currUserId, criteria)
		errx.PanicOnError(err)
		fmt.Println("\narticlesListDaf - favoritedBy:", articlePluses, "\n")
	}

	{
		err := daf.FavoriteDeleteDaf(ctx, tx, Favorites[1].ArticleId, Favorites[1].UserId)
		errx.PanicOnError(err)

		currUserId := users[0].Id

		criteria := rpc.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: util.PointerOf(users[2].Username),
			Limit:       nil,
			Offset:      nil,
		}
		articlePluses, err := daf.ArticlesListDaf(ctx, tx, currUserId, criteria)
		errx.PanicOnError(err)
		fmt.Println("\narticlesListDaf - favoritedBy:", articlePluses, "\n")
	}

	err = tx.Commit(ctx)
	errx.PanicOnError(err)
}
