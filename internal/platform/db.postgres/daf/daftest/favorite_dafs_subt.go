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
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx/dbpgxtest"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/pvillela/go-foa-realworld/internal/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

var favoriteDafsSubt = dbpgxtest.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {

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

		actual := testutil.ArticlePlusesToArticles(returned)

		expected0 := util.SliceFilter(mdb.ArticlePlusGetAll(favoritedBy),
			func(ap model.ArticlePlus) bool {
				return ap.Favorited
			})
		expected := testutil.ArticlePlusesToArticles(expected0)

		//_, _ = spew.Println("********* actual", actual)
		//_, _ = spew.Println("********* expected", expected)

		assert.ElementsMatch(t, expected, actual, msg)
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

		assert.ElementsMatch(t, expected, returned, msg)
	}

	{
		msg := "FavoriteDeleteDafI - delete existing favorite"

		currUsername := username1
		favoritedBy := username3
		slug := slug1

		articleId := mdb.ArticleGetBySlug(slug).Id
		userId := mdb.UserGetByName(favoritedBy).Id
		currUserId := mdb.UserGetByName(currUsername).Id

		err := daf.FavoriteDeleteDafI(ctx, tx, articleId, userId)
		errx.PanicOnError(err)

		mdb.FavoritedDelete(favoritedBy, slug)

		criteria := model.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: &favoritedBy,
			Limit:       nil,
			Offset:      nil,
		}
		returned, err := daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
		errx.PanicOnError(err)
		//fmt.Println("\narticlesListDaf - favoritedBy:\n", returned)

		actual := testutil.ArticlePlusesToArticles(returned)

		expected0 := util.SliceFilter(mdb.ArticlePlusGetAll(favoritedBy),
			func(ap model.ArticlePlus) bool {
				return ap.Favorited
			})
		expected := testutil.ArticlePlusesToArticles(expected0)

		assert.ElementsMatch(t, expected, actual, msg)
	}

	{
		msg := "FavoriteDeleteDafI - attempt to delete inexistent favorite"

		favoritedBy := username3
		slug := slug1

		articleId := mdb.ArticleGetBySlug(slug).Id
		userId := mdb.UserGetByName(favoritedBy).Id

		err := daf.FavoriteDeleteDafI(ctx, tx, articleId, userId)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg)
	}
})
