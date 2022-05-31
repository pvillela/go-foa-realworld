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
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/stretchr/testify/assert"
	"testing"
)

var articleDafsSubt = dbpgx.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {
	currUser, _ := mdb.UserGet("pvillela")
	author, _ := mdb.UserGet("joebloe")

	{
		criteria := model.ArticleCriteria{
			Tag:         nil,
			Author:      &author.Username,
			FavoritedBy: nil,
			Limit:       nil,
			Offset:      nil,
		}
		returned, err := daf.ArticlesListDafI(ctx, tx, currUser.Id, criteria)
		errx.PanicOnError(err)

		expected := util.SliceFilter(mdb.ArticlesPlus(), func(a model.ArticlePlus) bool {
			return a.Author.Username == author.Username
		})

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
		returned, err := daf.ArticlesListDafI(ctx, tx, currUser.Id, criteria)
		errx.PanicOnError(err)

		var expected []model.ArticlePlus

		//fmt.Println("\ncore info returned:", returned)
		//fmt.Println("\ncore info expected:", expected)

		assert.ElementsMatch(t, expected, returned)
	}

	{
		article := articles[1]

		returned, err := daf.ArticleGetBySlugDafI(ctx, tx, currUser.Id, article.Slug)
		errx.PanicOnError(err)

		expected := model.ArticlePlus_FromArticle(article, model.Profile_FromUser(&author, false))

		assert.Equal(t, expected, returned)
	}

	{
		pArticle := &articles[0]
		pArticle.Title = "A very interesting subject"
		err := daf.ArticleUpdateDafI(ctx, tx, pArticle)
		errx.PanicOnError(err)

		returned, err := daf.ArticleGetBySlugDafI(ctx, tx, currUser.Id, pArticle.Slug)
		errx.PanicOnError(err)

		expected := model.ArticlePlus_FromArticle(*pArticle, model.Profile_FromUser(&author, false))

		assert.Equal(t, expected, returned)
	}

	{
		returned, err := daf.ArticlesFeedDafI(ctx, tx, currUser.Id, nil, nil)
		errx.PanicOnError(err)

		var expected []model.ArticlePlus

		assert.ElementsMatch(t, expected, returned)
	}
})
