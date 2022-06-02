/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
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

	type favoriteSourceT struct {
		username string
		slug     string
	}

	favoriteSources := []favoriteSourceT{
		{
			username: username1,
			slug:     slug1,
		},
		{
			username: username3,
			slug:     slug1,
		},
	}

	for _, fsrc := range favoriteSources {
		username := fsrc.username
		slug := fsrc.slug
		articleId := mdb.ArticleGetBySlug(slug).Id
		userId := mdb.UserGetByName(username).Id
		err := daf.FavoriteCreateDafI(ctx, tx, articleId, userId)
		errx.PanicOnError(err)
		mdb.FavoritePut(username, slug)
	}

	// Tests

	{
		msg := "Get articles favorited by a given user who has favorites."

		currUsername := username1
		favoritedBy := username3
		slug := slug1

		currUser := mdb.UserGetByName(currUsername)

		criteria := model.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: &favoritedBy,
			Limit:       nil,
			Offset:      nil,
		}
		returned, err := daf.ArticlesListDafI(ctx, tx, currUser.Id, criteria)
		errx.PanicOnError(err)
		//fmt.Println("\narticlesListDaf - favoritedBy:", articlePluses, "\n")

		favoritedArticlePlus := mdb.ArticlePlusGet(currUsername, slug)

		expected := []model.ArticlePlus{favoritedArticlePlus}

		assert.Equal(t, expected, returned, msg)
	}

	{
		msg := "Get articles favorited by a given user who does not have favorites."

		currUsername := username1
		favoritedBy := username2

		currUser := mdb.UserGetByName(currUsername)

		criteria := model.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: &favoritedBy,
			Limit:       nil,
			Offset:      nil,
		}
		returned, err := daf.ArticlesListDafI(ctx, tx, currUser.Id, criteria)
		errx.PanicOnError(err)
		//fmt.Println("\narticlesListDaf - favoritedBy:", articlePluses, "\n")

		expected := []model.ArticlePlus{}

		assert.Equal(t, expected, returned, msg)
	}

	{
		msg := "FavoriteDeleteDafI - delete existing favorite"

		username := username1
		currUsername := username2
		slug := slug1

		articleId := mdb.ArticleGetBySlug(slug).Id
		userId := mdb.UserGetByName(username).Id
		currUserId := mdb.UserGetByName(currUsername).Id

		err := daf.FavoriteDeleteDafI(ctx, tx, articleId, userId)
		errx.PanicOnError(err)

		mdb.FavoritedDelete(username, slug)

		criteria := model.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: &username,
			Limit:       nil,
			Offset:      nil,
		}
		returned, err := daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
		errx.PanicOnError(err)
		//fmt.Println("\narticlesListDaf - favoritedBy:\n", returned)

		expected := util.SliceFilter(mdb.ArticlePlusGetAll(currUsername), func(ap model.ArticlePlus) bool {
			return ap.Favorited
		})

		assert.ElementsMatch(t, returned, expected, msg)
	}

	{
		msg := "FavoriteDeleteDafI - attempt to delete inexistent favorite"

		username := username1
		slug := slug1

		articleId := mdb.ArticleGetBySlug(slug).Id
		userId := mdb.UserGetByName(username).Id

		err := daf.FavoriteDeleteDafI(ctx, tx, articleId, userId)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg)
	}
})
