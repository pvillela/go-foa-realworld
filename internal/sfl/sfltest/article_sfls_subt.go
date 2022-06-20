/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfltest

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
	"github.com/pvillela/go-foa-realworld/internal/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

///////////////////
// Shared constants and data

const (
	slug1 = "anintsubj"
	slug2 = "adullsubj"
)

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

///////////////////
// Tests

func articleCreateSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	articleCreateSfl := sfl.ArticleCreateSflC(db)

	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	assert.NoError(t, err)

	{
		msg := "article_create_sfl - valid article"

		for _, aa := range authorsAndArticles {
			authorname := aa.Authorname
			reqCtx := web.RequestContext{
				Username: authorname,
				Token:    &jwt.Token{},
			}
			article := aa.Article

			in := rpc.ArticleCreateIn{Article: rpc.ArticleCreateIn0{
				Title:       article.Title,
				Description: article.Description,
				Body:        article.Body,
				TagList:     nil,
			}}

			articleOut, err := articleCreateSfl(ctx, reqCtx, in)
			assert.NoError(t, err)

			user, err := testutil.UserGetByName(db, ctx, authorname)
			assert.NoError(t, err)

			expected := rpc.ArticleOut{Article: model.ArticlePlus{
				Id:             articleOut.Article.Id, // not independently checked
				Slug:           util.Slug(article.Title),
				Author:         model.Profile_FromUser(user, false),
				Title:          article.Title,
				Description:    article.Description,
				Body:           article.Body,
				TagList:        article.TagList,
				CreatedAt:      articleOut.Article.CreatedAt, // not independently checked
				UpdatedAt:      articleOut.Article.UpdatedAt, // not independently checked
				Favorited:      false,
				FavoritesCount: 0,
			}}

			assert.Equal(t, expected, articleOut, msg+" - output must align with input")
		}
	}
}
