/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf/daftest/memdb"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

var mdb memdb.MDb

var users = []model.User{
	{
		Username:     "pvillela",
		Email:        "foo@bar.com",
		PasswordHash: "dakfljads0fj",
		PasswordSalt: "2af8d0b50a",
		Bio:          util.PointerFromValue("I am me."),
		ImageLink:    nil,
	},
	{
		Username:     "joebloe",
		Email:        "joe@bloe.com",
		PasswordHash: "9zdakfljads0",
		PasswordSalt: "3ba9e9c611",
		Bio:          util.PointerFromValue("Famous person."),
		ImageLink:    util.PointerFromValue("https://myimage.com"),
	},
	{
		Username:     "johndoe",
		Email:        "johndoe@foo.com",
		PasswordHash: "09fs8asfoasi",
		PasswordSalt: "0000000000",
		Bio:          util.PointerFromValue("Average guy."),
		ImageLink:    util.PointerFromValue("https://johndooeimage.com"),
	},
}

type AuthorAndArticle struct {
	Authorname string
	Article    model.Article
}

var authorsAndArticles = []AuthorAndArticle{
	{
		Authorname: users[1].Username,
		Article: model.Article{
			Title:       "An interesting subject",
			Slug:        "anintsubj",
			Description: "Story about an interesting subject.",
			Body:        util.PointerFromValue("I met this interesting subject a long time ago."),
		},
	},
	{
		Authorname: users[1].Username,
		Article: model.Article{
			Title:       "A dull story",
			Slug:        "adullsubj",
			Description: "Narrative about something dull.",
			Body:        util.PointerFromValue("This is so dull, bla, bla, bla."),
		},
	},
}

func setupData(ctx context.Context, tx pgx.Tx) {
	for i, _ := range users {
		user := users[i]
		recCtx, err := daf.UserCreateExplicitTxDafI(ctx, tx, &user)
		errx.PanicOnError(err)
		//_, _ = spew.Printf("user from Create: %v\n", user)
		logrus.Debug("user from Create:", user)
		logrus.Debug("recCtx from Create:", recCtx)

		mdb.UserUpsert(user.Username, user, recCtx)
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
