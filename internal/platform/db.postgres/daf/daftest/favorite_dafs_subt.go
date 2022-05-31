/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/stretchr/testify/assert"
	"testing"
)

var favoriteDafsSubt = dbpgx.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {
	// Create favorites.

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
		err := daf.FavoriteCreateDafI(ctx, tx, Favorites[i].ArticleId, Favorites[i].UserId)
		errx.PanicOnError(err)
	}

	// Tests

	authors := authors()

	{
		msg := "Get articles favorited by a given user."

		currUserId := users[0].Id
		favoritedBy := users[2]

		criteria := model.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: &favoritedBy.Username,
			Limit:       nil,
			Offset:      nil,
		}
		returned, err := daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
		errx.PanicOnError(err)
		//fmt.Println("\narticlesListDaf - favoritedBy:", articlePluses, "\n")

		favoritedArticle := articles[0]
		favoritedArticlePlus := model.ArticlePlus_FromArticle(
			favoritedArticle,
			model.Profile_FromUser(&authors[0], false),
		)
		expected := []model.ArticlePlus{favoritedArticlePlus}

		assert.Equal(t, expected, returned, msg)
	}

	{
		currUserId := users[2].Id

		criteria := model.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: util.PointerFromValue(users[0].Username),
			Limit:       nil,
			Offset:      nil,
		}
		articlePluses, err := daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
		errx.PanicOnError(err)
		fmt.Println("\narticlesListDaf - favoritedBy:\n", articlePluses)
	}

	{
		err := daf.FavoriteDeleteDafI(ctx, tx, Favorites[1].ArticleId, Favorites[1].UserId)
		errx.PanicOnError(err)

		currUserId := users[0].Id

		criteria := model.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: util.PointerFromValue(users[2].Username),
			Limit:       nil,
			Offset:      nil,
		}
		articlePluses, err := daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
		errx.PanicOnError(err)
		fmt.Println("\narticlesListDaf - favoritedBy:\n", articlePluses)
	}
})
