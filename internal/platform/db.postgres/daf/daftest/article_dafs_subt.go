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
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupArticles(ctx context.Context, tx pgx.Tx) {
	type AuthorAndArticle struct {
		Authorname string
		Article    model.Article
	}

	var authorsAndArticles = []AuthorAndArticle{
		{
			Authorname: "joebloe",
			Article: model.Article{
				Title:       "An interesting subject",
				Slug:        "anintsubj",
				Description: "Story about an interesting subject.",
				Body:        util.PointerFromValue("I met this interesting subject a long time ago."),
			},
		},
		{
			Authorname: "joebloe",
			Article: model.Article{
				Title:       "A dull story",
				Slug:        "adullsubj",
				Description: "Narrative about something dull.",
				Body:        util.PointerFromValue("This is so dull, bla, bla, bla."),
			},
		},
	}

	for i := range authorsAndArticles {
		authorname := authorsAndArticles[i].Authorname
		author, _ := mdb.UserGet(authorname)
		if author == (model.User{}) {
			panic(fmt.Sprintf("invalid username %v", authorname))
		}
		article := authorsAndArticles[i].Article
		article.AuthorId = author.Id
		err := daf.ArticleCreateDafI(ctx, tx, &article)
		errx.PanicOnError(err)
		logrus.Debug("article from Create:", article)

		mdb.ArticlePlusUpsert(article, false, author, false)
	}
}

var articleDafsSubt = dbpgx.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {
	
	{
		setupArticles(ctx, tx)
	}

	currUser, _ := mdb.UserGet("pvillela")
	author, _ := mdb.UserGet("joebloe")

	{
		msg := "ArticlesListDafI - select author"

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

		assert.ElementsMatch(t, expected, returned, msg)
	}

	{
		msg := "ArticlesListDafI - select tag"

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

		assert.ElementsMatch(t, expected, returned, msg)
	}

	{
		msg := "ArticleGetBySlugDafI"

		slug := "adullsubj"

		returned, err := daf.ArticleGetBySlugDafI(ctx, tx, currUser.Id, slug)
		errx.PanicOnError(err)

		expected := mdb.ArticlePlusGet(slug)

		assert.Equal(t, expected, returned, msg)
	}

	{
		msg := "ArticlesFeedDafI"

		returned, err := daf.ArticlesFeedDafI(ctx, tx, currUser.Id, nil, nil)
		errx.PanicOnError(err)

		var expected []model.ArticlePlus

		assert.ElementsMatch(t, expected, returned, msg)
	}
})
