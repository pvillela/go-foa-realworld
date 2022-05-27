/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func articlePlusToCore(_ int, ap model.ArticlePlus) model.ArticlePlus {
	return model.ArticlePlus{
		Slug:        ap.Slug,
		Author:      ap.Author,
		Title:       ap.Title,
		Description: ap.Description,
		Body:        ap.Body,
		TagList:     ap.TagList,
	}
}

func articleToCore(authors []model.User, follows bool) func(i int, a model.Article) model.ArticlePlus {
	return func(i int, a model.Article) model.ArticlePlus {
		return model.ArticlePlus{
			Slug:        a.Slug,
			Author:      model.Profile_FromUser(&authors[i], follows),
			Title:       a.Title,
			Description: a.Description,
			Body:        a.Body,
			TagList:     a.TagList,
		}
	}
}

func dbTestWithTransaction(db dbpgx.Db, ctx context.Context, block func(tx pgx.Tx)) {
	block1 := func(ctx context.Context, tx pgx.Tx) (types.Unit, error) {
		block(tx)
		return types.UnitV, nil
	}
	_, err := dbpgx.Db_WithTransaction(db, ctx, block1)
	errx.PanicOnError(err)
}

type DafSubt func(
	ctx context.Context,
	tx pgx.Tx,
	t *testing.T,
)

func dbTestWithTransactionL(
	db dbpgx.Db,
	f DafSubt,
) func(ctx context.Context, t *testing.T) {
	return util.LiftContextualizer1V(dbpgx.Db_WithTransaction[types.Unit], db, f)
}

func articleDafsSubt0(db dbpgx.Db, ctx context.Context, t *testing.T) {
	dbTestWithTransaction(db, ctx, func(tx pgx.Tx) {
		currUser := users[0]
		authors := []model.User{users[1], users[1]}
		author := authors[0]

		{
			criteria := model.ArticleCriteria{
				Tag:         nil,
				Author:      &author.Username,
				FavoritedBy: nil,
				Limit:       nil,
				Offset:      nil,
			}
			articlePluses, err := daf.ArticlesListDafI(ctx, tx, currUser.Id, criteria)
			errx.PanicOnError(err)

			returned := util.SliceMap(articlePluses, articlePlusToCore)
			expected := util.SliceMap(articles, articleToCore(authors, false))

			//fmt.Println("\ncoreInfoReturned:", coreInfoReturned)
			//fmt.Println("\ncoreInfoExpected:", coreInfoExpected)

			assert.ElementsMatch(t, expected, returned)
		}

		{
			criteria := model.ArticleCriteria{
				Tag:         util.PointerFromValue("FOOTAG"),
				Author:      nil,
				FavoritedBy: nil,
				Limit:       nil,
				Offset:      nil,
			}
			articlePluses, err := daf.ArticlesListDafI(ctx, tx, currUser.Id, criteria)
			errx.PanicOnError(err)

			returned := util.SliceMap(articlePluses, articlePlusToCore)
			var expected []model.ArticlePlus

			//fmt.Println("\ncore info returned:", returned)
			//fmt.Println("\ncore info expected:", expected)

			assert.ElementsMatch(t, expected, returned)
		}

		{
			articleFromDb, err := daf.ArticleGetBySlugDafI(ctx, tx, currUser.Id, articles[1].Slug)
			errx.PanicOnError(err)

			returned := articlePlusToCore(-1, articleFromDb)
			expected := articleToCore(authors, false)(0, articles[1])

			//_, _ = spew.Println("\nArticleGetBySlugDaf:", articleFromDb)

			assert.Equal(t, expected, returned)
		}

		{
			pArticle := &articles[0]
			pArticle.Title = "A very interesting subject"
			err := daf.ArticleUpdateDafI(ctx, tx, pArticle)
			errx.PanicOnError(err)

			articleFromDb, err := daf.ArticleGetBySlugDafI(ctx, tx, currUser.Id, pArticle.Slug)
			errx.PanicOnError(err)

			returned := articlePlusToCore(-1, articleFromDb)
			expected := articleToCore(authors, false)(0, *pArticle)

			//_, _ = spew.Println("\nAfter update:", articleFromDb)

			assert.Equal(t, expected, returned)
		}

		{
			articlePluses, err := daf.ArticlesFeedDafI(ctx, tx, currUser.Id, nil, nil)
			errx.PanicOnError(err)

			returned := util.SliceMap(articlePluses, articlePlusToCore)
			var expected []model.ArticlePlus

			//_, _ = spew.Println("\nArticlesFeedDaf returned:", returned)
			//_, _ = spew.Println("\nArticlesFeedDaf expected:", expected)

			assert.ElementsMatch(t, expected, returned)
		}
	})
}

func articleDafsSubt1(ctx context.Context, tx pgx.Tx, t *testing.T) {
	currUser := users[0]
	authors := []model.User{users[1], users[1]}
	author := authors[0]

	{
		criteria := model.ArticleCriteria{
			Tag:         nil,
			Author:      &author.Username,
			FavoritedBy: nil,
			Limit:       nil,
			Offset:      nil,
		}
		articlePluses, err := daf.ArticlesListDafI(ctx, tx, currUser.Id, criteria)
		errx.PanicOnError(err)

		returned := util.SliceMap(articlePluses, articlePlusToCore)
		expected := util.SliceMap(articles, articleToCore(authors, false))

		//fmt.Println("\ncoreInfoReturned:", coreInfoReturned)
		//fmt.Println("\ncoreInfoExpected:", coreInfoExpected)

		assert.ElementsMatch(t, expected, returned)
	}

	{
		criteria := model.ArticleCriteria{
			Tag:         util.PointerFromValue("FOOTAG"),
			Author:      nil,
			FavoritedBy: nil,
			Limit:       nil,
			Offset:      nil,
		}
		articlePluses, err := daf.ArticlesListDafI(ctx, tx, currUser.Id, criteria)
		errx.PanicOnError(err)

		returned := util.SliceMap(articlePluses, articlePlusToCore)
		var expected []model.ArticlePlus

		//fmt.Println("\ncore info returned:", returned)
		//fmt.Println("\ncore info expected:", expected)

		assert.ElementsMatch(t, expected, returned)
	}

	{
		articleFromDb, err := daf.ArticleGetBySlugDafI(ctx, tx, currUser.Id, articles[1].Slug)
		errx.PanicOnError(err)

		returned := articlePlusToCore(-1, articleFromDb)
		expected := articleToCore(authors, false)(0, articles[1])

		//_, _ = spew.Println("\nArticleGetBySlugDaf:", articleFromDb)

		assert.Equal(t, expected, returned)
	}

	{
		pArticle := &articles[0]
		pArticle.Title = "A very interesting subject"
		err := daf.ArticleUpdateDafI(ctx, tx, pArticle)
		errx.PanicOnError(err)

		articleFromDb, err := daf.ArticleGetBySlugDafI(ctx, tx, currUser.Id, pArticle.Slug)
		errx.PanicOnError(err)

		returned := articlePlusToCore(-1, articleFromDb)
		expected := articleToCore(authors, false)(0, *pArticle)

		//_, _ = spew.Println("\nAfter update:", articleFromDb)

		assert.Equal(t, expected, returned)
	}

	{
		articlePluses, err := daf.ArticlesFeedDafI(ctx, tx, currUser.Id, nil, nil)
		errx.PanicOnError(err)

		returned := util.SliceMap(articlePluses, articlePlusToCore)
		var expected []model.ArticlePlus

		//_, _ = spew.Println("\nArticlesFeedDaf returned:", returned)
		//_, _ = spew.Println("\nArticlesFeedDaf expected:", expected)

		assert.ElementsMatch(t, expected, returned)
	}
}
