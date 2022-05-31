/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/sirupsen/logrus"
)

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

var recCtxUsers = make([]daf.RecCtxUser, len(users))

var articles = []model.Article{
	{
		Title:       "An interesting subject",
		Slug:        "anintsubj",
		Description: "Story about an interesting subject.",
		Body:        util.PointerFromValue("I met this interesting subject a long time ago."),
	},
	{
		Title:       "A dull story",
		Slug:        "adullsubj",
		Description: "Narrative about something dull.",
		Body:        util.PointerFromValue("This is so dull, bla, bla, bla."),
	},
}

func authors() []model.User {
	authors := []model.User{users[1], users[1]}

	if len(authors) != len(articles) {
		panic("len(authors()) != len(articles)")
	}
	return authors
}

func setupData(ctx context.Context, tx pgx.Tx) {
	for i, _ := range users {
		recCtx, err := daf.UserCreateExplicitTxDafI(ctx, tx, &users[i])
		errx.PanicOnError(err)
		recCtxUsers[i] = recCtx
		_, _ = spew.Printf("user from Create: %v\n", users[i])
		logrus.Debug("recCtx from Create:", recCtx)
	}

	authors := authors()

	for i, _ := range articles {
		articles[i].AuthorId = authors[i].Id
		err := daf.ArticleCreateDafI(ctx, tx, &articles[i])
		errx.PanicOnError(err)
		logrus.Debug("article from Create:", articles[i])
	}
}

type MDb struct {
	Users        MUsers
	RecCtxUsers  MRecCtxUsers
	ArticlesPlus MArticlesPlus
	Favorites    MFavorites
	Followings   MFollowings
	Comments     MComments
	Tags         MTags
}

// key is Slug
type MArticlesPlus map[string]model.ArticlePlus

// key is Username
type MUsers map[string]model.User

// key is Username
type MRecCtxUsers map[string]daf.RecCtxUser

// key is Username, value is Slug
type MFavorites map[string]string

type FollowingsKey struct {
	FollowerName string
	FolloweeName string
}

type MFollowings map[FollowingsKey]model.Following

type MComment struct {
	Username string
	Slug     string
	Comment  model.Comment
}

type MComments []MComment

type MTags map[string]types.Unit

var mdb MDb

type AuthorAndArticle struct {
	AuthorName string
	Article    model.Article
}

var authorsAndArticles = []AuthorAndArticle{
	{users[1].Username, articles[0]},
	{users[1].Username, articles[1]},
}

func setupData1(ctx context.Context, tx pgx.Tx) {
	for i, _ := range users {
		user := users[i]
		recCtx, err := daf.UserCreateExplicitTxDafI(ctx, tx, &user)
		errx.PanicOnError(err)
		mdb.Users[user.Username] = user
		mdb.RecCtxUsers[user.Username] = recCtx
		_, _ = spew.Printf("user from Create: %v\n", user)
		logrus.Debug("recCtx from Create:", recCtx)
	}

	for i := range authorsAndArticles {
		authorname := authorsAndArticles[i].AuthorName
		author := mdb.Users[authorname]
		if author == (model.User{}) {
			panic(fmt.Sprintf("invalid username %v", authorname))
		}
		article := authorsAndArticles[i].Article
		article.AuthorId = author.Id
		err := daf.ArticleCreateDafI(ctx, tx, &article)
		errx.PanicOnError(err)
		mdb.ArticlesPlus[article.Slug] = model.ArticlePlus_FromArticle(
			article,
			false,
			model.Profile_FromUser(&author, false),
		)
		logrus.Debug("article from Create:", articles[i])
	}
}
