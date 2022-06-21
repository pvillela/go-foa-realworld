/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx/dbpgxtest"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	slug1 = "anintsubj"
	slug2 = "adullsubj"
)

func setupArticles(ctx context.Context, tx pgx.Tx) {
	type AuthorAndArticle struct {
		Authorname string
		Article    model.Article
	}

	var authorsAndArticles = []AuthorAndArticle{
		{
			Authorname: username2,
			Article: model.Article{
				Title:       "An interesting subject",
				Slug:        slug1,
				Description: "Story about an interesting subject.",
				Body:        util.PointerFromValue("I met this interesting subject a long time ago."),
			},
		},
		{
			Authorname: username2,
			Article: model.Article{
				Title:       "A dull story",
				Slug:        slug2,
				Description: "Narrative about something dull.",
				Body:        util.PointerFromValue("This is so dull, bla, bla, bla."),
			},
		},
	}

	for i := range authorsAndArticles {
		authorname := authorsAndArticles[i].Authorname
		author := mdb.UserGetByName(authorname)
		if author == (model.User{}) {
			panic(fmt.Sprintf("invalid username %v", authorname))
		}
		article := authorsAndArticles[i].Article
		article.AuthorId = author.Id
		err := daf.ArticleCreateDaf(ctx, tx, &article)
		errx.PanicOnError(err)
		logrus.Debug("article from Create:", article)

		mdb.ArticleUpsert(article)
	}
}

var articleDafsSubt = dbpgxtest.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {

	{
		setupArticles(ctx, tx)
	}

	currUser := mdb.UserGetByName(username1)
	author := mdb.UserGetByName(username2)

	{
		msg := "ArticleUpdateDaf"

		slug := slug1

		existingArticle := mdb.ArticleGetBySlug(slug)
		changedArticle := existingArticle
		newBody := util.PointerFromValue(*existingArticle.Body + " And so on and on.")
		changedArticle.Body = newBody

		err := daf.ArticleUpdateDaf(ctx, tx, &changedArticle)
		assert.NoError(t, err)

		assert.Equal(t, newBody, changedArticle.Body, msg+" - changedArticle.Body must equal newBody")
		assert.Greater(t, changedArticle.UpdatedAt, existingArticle.UpdatedAt,
			msg+" - changedArticle must have an UpdatedAt value greater than existingArticle")

		changedArticlePrime := changedArticle
		changedArticlePrime.UpdatedAt = existingArticle.UpdatedAt
		changedArticlePrime.Body = existingArticle.Body

		assert.Equal(t, existingArticle, changedArticlePrime,
			msg+" - changed article must equal existing article except for Body and ChangedAt")

		returned, err := daf.ArticleGetBySlugDaf(ctx, tx, currUser.Id, slug)
		assert.NoError(t, err)

		//criteria := model.ArticleCriteria{
		//	Tag:         nil,
		//	Author:      &author.Username,
		//	FavoritedBy: nil,
		//	Limit:       nil,
		//	Offset:      nil,
		//}
		//aps, err := daf.ArticlesListDaf(ctx, tx, currUser.Id, criteria)
		//assert.NoError(t, err)
		//spew.Println("************* aps", aps)

		actual := returned.ToArticle()
		expected := changedArticle

		//spew.Println("************* returned", returned)
		//fmt.Println("************* expected", expected)
		//fmt.Println("************* actual", actual)

		assert.Equal(t, expected, actual, msg+" - retrieved must equal in-memory")

		mdb.ArticleUpsert(changedArticle)
	}

	{
		msg := "ArticlesListDaf - select author"

		criteria := model.ArticleCriteria{
			Tag:         nil,
			Author:      &author.Username,
			FavoritedBy: nil,
			Limit:       nil,
			Offset:      nil,
		}
		returned, err := daf.ArticlesListDaf(ctx, tx, currUser.Id, criteria)
		assert.NoError(t, err)

		expected := util.SliceFilter(mdb.ArticlePlusGetAll(author.Username), func(a model.ArticlePlus) bool {
			return a.Author.Username == author.Username
		})

		assert.ElementsMatch(t, expected, returned, msg)
	}

	{
		msg := "ArticlesListDaf - select tag"

		criteria := model.ArticleCriteria{
			Tag:         util.PointerFromValue("FOOTAG"),
			Author:      nil,
			FavoritedBy: nil,
			Limit:       nil,
			Offset:      nil,
		}
		returned, err := daf.ArticlesListDaf(ctx, tx, currUser.Id, criteria)
		assert.NoError(t, err)

		var expected []model.ArticlePlus

		//fmt.Println("\ncore info returned:", returned)
		//fmt.Println("\ncore info expected:", expected)

		assert.ElementsMatch(t, expected, returned, msg)
	}

	{
		msg := "ArticleGetBySlugDaf"

		slug := slug2

		returned, err := daf.ArticleGetBySlugDaf(ctx, tx, currUser.Id, slug)
		assert.NoError(t, err)

		expected := mdb.ArticlePlusGet(currUser.Username, slug)

		assert.Equal(t, expected, returned, msg)
	}

	{
		msg := "ArticlesFeedDaf"

		returned, err := daf.ArticlesFeedDaf(ctx, tx, currUser.Id, nil, nil)
		assert.NoError(t, err)

		var expected []model.ArticlePlus

		assert.ElementsMatch(t, expected, returned, msg)
	}
})
