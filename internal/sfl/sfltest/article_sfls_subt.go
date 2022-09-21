/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfltest

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
	"github.com/pvillela/go-foa-realworld/internal/sfl/boot"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/testutil"
	"github.com/stretchr/testify/assert"
)

///////////////////
// Shared constants and data

type AuthorAndArticle struct {
	Authorname string
	Article    model.Article
}

var authorsAndArticles = []AuthorAndArticle{
	{
		Authorname: username2,
		Article: model.Article{
			Title:       "An interesting subject",
			Description: "Story about an interesting subject.",
			Body:        util.PointerOf("I met this interesting subject a long time ago."),
		},
	},
	{
		Authorname: username2,
		Article: model.Article{
			Title:       "A dull story",
			Description: "Narrative about something dull.",
			Body:        util.PointerOf("This is so dull, bla, bla, bla."),
		},
	},
	{
		Authorname: username2,
		Article: model.Article{
			Title:       "An article to be deleted",
			Description: "Stuff about an article to be deleted.",
			Body:        util.PointerOf("This is an article to be deleted, bla, bla, bla."),
		},
	},
}

///////////////////
// Tests

func articleCreateSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	boot.ArticleCreateSflCfgAdapter = TestCfgAdapterOf(db)
	articleCreateSfl := boot.ArticleCreateSflBoot(nil)

	{
		msg := "article_create_sfl - valid article"

		for _, aa := range util.SliceReverse(authorsAndArticles) {
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

	{
		msg := "article_create_sfl - existing title which implies existing slug"

		// Try to recreate articles with existing slugs
		for _, aa := range authorsAndArticles {
			authorname := aa.Authorname
			reqCtx := web.RequestContext{
				Username: authorname,
				Token:    &jwt.Token{},
			}
			article := aa.Article

			in := rpc.ArticleCreateIn{Article: rpc.ArticleCreateIn0{
				Title:       article.Title,
				Description: "dummy description",
				Body:        util.PointerOf("dummy body"),
				TagList:     nil,
			}}

			_, err := articleCreateSfl(ctx, reqCtx, in)
			returnedErrxKind := dbpgx.ClassifyError(err)
			expectedErrxKind := dbpgx.DbErrUniqueViolation
			expectedErrMsgPrefix := "DbErrUniqueViolation[duplicate article slug"

			assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when username or email is not unique")
			assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username or email is not unique")
		}
	}
}

func articleDeleteSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	boot.ArticleDeleteSflCfgAdapter = TestCfgAdapterOf(db)
	articleDeleteSfl := boot.ArticleDeleteSflC(nil)

	{
		msg := "article_delete_sfl - existing article authored by current user"

		aa := authorsAndArticles[2]
		article := aa.Article
		slug := util.Slug(article.Title)
		authorname := aa.Authorname

		reqCtx := web.RequestContext{
			Username: authorname,
			Token:    &jwt.Token{},
		}

		_, err := articleDeleteSfl(ctx, reqCtx, slug)
		assert.NoError(t, err, msg)

		_, err = testutil.ArticleGetBySlug(db, ctx, authorname, slug)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[article slug"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - retrieval of deleted article must fail with appropriate error kind when username or email is not unique")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - retrieval of deleted article must fail with appropriate error message when username or email is not unique")
	}

	{
		msg := "article_delete_sfl - inexistenet article"

		aa := authorsAndArticles[2]
		article := aa.Article
		slug := util.Slug(article.Title)
		authorname := aa.Authorname

		reqCtx := web.RequestContext{
			Username: authorname,
			Token:    &jwt.Token{},
		}

		_, err := articleDeleteSfl(ctx, reqCtx, slug)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[article slug"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when username or email is not unique")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username or email is not unique")
	}

	{
		msg := "article_delete_sfl - attempted by non-author"

		aa := authorsAndArticles[0]
		article := aa.Article
		slug := util.Slug(article.Title)
		authorname := username1

		reqCtx := web.RequestContext{
			Username: authorname,
			Token:    &jwt.Token{},
		}

		_, err := articleDeleteSfl(ctx, reqCtx, slug)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := bf.ErrUnauthorizedUser
		expectedErrMsgPrefix := "ErrUnauthorizedUser[user"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when username or email is not unique")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username or email is not unique")
	}
}

func articleFavoriteSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	boot.ArticleFavoriteSflCfgAdapter = TestCfgAdapterOf(db)
	articleFavoriteSfl := boot.ArticleFavoriteSflC(nil)

	{
		msg := "article_favorite_sfl - existing article, not yet favorited"

		currUsername := username1
		aa := authorsAndArticles[0]
		article := aa.Article
		slug := util.Slug(article.Title)

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		out, err := articleFavoriteSfl(ctx, reqCtx, slug)
		assert.NoError(t, err, msg)

		assert.True(t, out.Article.Favorited, msg+" - Favorited attribute of output must be true")
		assert.Equal(t, article.Description, out.Article.Description, msg+" - Description attribute must not change")
		assert.Equal(t, article.Body, out.Article.Body, msg+" - Body attribute must not change")
		assert.Equal(t, article.FavoritesCount+1, out.Article.FavoritesCount, msg+" - FavoritesCount attribute must be incremented")
	}

	{
		msg := "article_favorite_sfl - article already favorited"

		currUsername := username1
		aa := authorsAndArticles[0]
		article := aa.Article
		slug := util.Slug(article.Title)

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		_, err := articleFavoriteSfl(ctx, reqCtx, slug)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrUniqueViolation
		expectedErrMsgPrefix := "DbErrUniqueViolation[article with ID"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when favoriting an already favorited article")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when favoriting an already favorited article")
	}

	{
		msg := "article_favorite_sfl - inexistent article"

		currUsername := username1
		slug := "dkdkddkd"

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		_, err := articleFavoriteSfl(ctx, reqCtx, slug)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[article slug"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when favoriting an inexistent article")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when favoriting an inexistent article")
	}
}

func articleGetSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	boot.ArticleGetSflCfgAdapter = TestCfgAdapterOf(db)
	articleGetSfl := boot.ArticleGetSflC(nil)

	{
		msg := "article_get_sfl - existing article"

		currUsername := username1
		aa := authorsAndArticles[0]
		article := aa.Article
		slug := util.Slug(article.Title)

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		out, err := articleGetSfl(ctx, reqCtx, slug)
		assert.NoError(t, err, msg)

		assert.True(t, out.Article.Favorited, msg+" - Favorited attribute of output must be true because it already was before")
		assert.Equal(t, article.Title, out.Article.Title, msg+" - Title attribute must not change")
		assert.Equal(t, article.Description, out.Article.Description, msg+" - Description attribute must not change")
		assert.Equal(t, article.Body, out.Article.Body, msg+" - Body attribute must not change")
	}

	{
		msg := "article_get_sfl - inexistent article"

		currUsername := username1
		slug := "dkdkddkd"

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		_, err := articleGetSfl(ctx, reqCtx, slug)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[article slug"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when favoriting an inexistent article")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when favoriting an inexistent article")
	}
}

func articleUnfavoriteSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	boot.ArticleUnfavoriteSflCfgAdapter = TestCfgAdapterOf(db)
	articleUnfavoriteSfl := boot.ArticleUnfavoriteSflC(nil)

	{
		msg := "article_unfavorite_sfl - existing article, previously favorited"

		currUsername := username1
		aa := authorsAndArticles[0]
		article := aa.Article
		slug := util.Slug(article.Title)

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		out, err := articleUnfavoriteSfl(ctx, reqCtx, slug)
		assert.NoError(t, err, msg)

		assert.False(t, out.Article.Favorited, msg+" - Favorited attribute of output must be false")
		assert.Equal(t, article.Description, out.Article.Description, msg+" - Description attribute must not change")
		assert.Equal(t, article.Body, out.Article.Body, msg+" - Body attribute must not change")
		assert.Equal(t, article.FavoritesCount, out.Article.FavoritesCount, msg+" - FavoritesCount attribute must go back to what it was initially")
	}

	{
		msg := "article_unfavorite_sfl - existing article, not previously favorited"

		currUsername := username1
		aa := authorsAndArticles[0]
		article := aa.Article
		slug := util.Slug(article.Title)

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		_, err := articleUnfavoriteSfl(ctx, reqCtx, slug)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[article with ID"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when unfavoriting an article not previously favorited")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when unfavoriting an article not previously favorited")
	}

	{
		msg := "article_unfavorite_sfl - inexistent article"

		currUsername := username1
		slug := "dkdkddkd"

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		_, err := articleUnfavoriteSfl(ctx, reqCtx, slug)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[article slug"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when unfavoriting an inexistent article")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when unfavoriting an inexistent article")
	}
}

func articleUpdateSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	boot.ArticleUpdateSflCfgAdapter = TestCfgAdapterOf(db)
	articleUpdateSfl := boot.ArticleUpdateSflBoot(nil)

	{
		msg := "article_update_sfl - existing article authored by current user"

		aa := authorsAndArticles[1]
		article := aa.Article
		slug := util.Slug(article.Title)
		authorname := aa.Authorname

		reqCtx := web.RequestContext{
			Username: authorname,
			Token:    &jwt.Token{},
		}

		in := rpc.ArticleUpdateIn{
			Article: rpc.ArticleUpdateIn0{
				Slug:        slug,
				Title:       util.PointerOf("New Title"),
				Description: util.PointerOf("New description"),
				Body:        util.PointerOf("New body"),
			},
		}

		_, err := articleUpdateSfl(ctx, reqCtx, in)
		assert.NoError(t, err, msg)

		retrieved, err := testutil.ArticleGetBySlug(db, ctx, authorname, slug)
		assert.NoError(t, err, msg)

		assert.Equal(t, *in.Article.Title, retrieved.Title, msg+" - Title attribute of retrieved article must match update")
		assert.Equal(t, *in.Article.Description, retrieved.Description, msg+" - Description attribute of retrieved article must match update")
		assert.Equal(t, in.Article.Body, retrieved.Body, msg+" - Body attribute of retrieved article must match update")
	}

	{
		msg := "article_update_sfl - inexistenet article"

		aa := authorsAndArticles[1]
		slug := "xxx"
		authorname := aa.Authorname

		reqCtx := web.RequestContext{
			Username: authorname,
			Token:    &jwt.Token{},
		}

		in := rpc.ArticleUpdateIn{
			Article: rpc.ArticleUpdateIn0{
				Slug:        slug,
				Title:       util.PointerOf("New Title"),
				Description: util.PointerOf("New description"),
				Body:        util.PointerOf("New body"),
			},
		}

		_, err := articleUpdateSfl(ctx, reqCtx, in)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[article slug"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when article doesn't exist")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when article doesn't exist")
	}

	{
		msg := "article_update_sfl - attempted by non-author"

		aa := authorsAndArticles[0]
		article := aa.Article
		slug := util.Slug(article.Title)

		authorname := username1
		reqCtx := web.RequestContext{
			Username: authorname,
			Token:    &jwt.Token{},
		}

		in := rpc.ArticleUpdateIn{
			Article: rpc.ArticleUpdateIn0{
				Slug:        slug,
				Title:       util.PointerOf("New Title"),
				Description: util.PointerOf("New description"),
				Body:        util.PointerOf("New body"),
			},
		}

		_, err := articleUpdateSfl(ctx, reqCtx, in)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := bf.ErrUnauthorizedUser
		expectedErrMsgPrefix := "ErrUnauthorizedUser[user"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when update attempt by non-author")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when update attempt by non-author")
	}
}

func articlesFeedSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	boot.ArticlesFeedSflCfgAdapter = TestCfgAdapterOf(db)
	articlesFeedSfl := boot.ArticlesFeedSflBoot(nil)

	{
		msg := "articles_feed_sfl - without limit and offset"

		currUsername := username1

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		in := rpc.ArticlesFeedIn{}

		out, err := articlesFeedSfl(ctx, reqCtx, in)
		assert.NoError(t, err, msg)
		assert.Equal(t, 2, len(out.Articles), msg+" - expected list with 2 articles")
	}

	{
		msg := "articles_feed_sfl - with limit and offset"

		currUsername := username1

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		in := rpc.ArticlesFeedIn{
			Limit:  util.PointerOf(1),
			Offset: util.PointerOf(1),
		}

		out, err := articlesFeedSfl(ctx, reqCtx, in)
		assert.NoError(t, err, msg)
		assert.Equal(t, 1, len(out.Articles), msg+" - expected list with one article")

		feedArticlce := out.Articles[0].Article.ToArticle()

		allArticles, err := testutil.ArticlesList(db, ctx, currUsername)
		assert.NoError(t, err, msg)
		assert.Equal(t, 2, len(allArticles), msg+" - retrieving all articles, expected list with 2 article")
		expectedArticle := allArticles[1]
		assert.Equal(t, expectedArticle, feedArticlce, msg+" - article from SFL must match second item in the list of all articles")
	}
}

func articlesListSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	boot.ArticlesListSflCfgAdapter = TestCfgAdapterOf(db)
	articlesListSfl := boot.ArticlesListSflBoot(nil)

	{
		msg := "articles_list_sfl - without limit and offset"

		currUsername := username1

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		in := rpc.ArticleCriteria{
			Author: util.PointerOf(username2),
		}

		out, err := articlesListSfl(ctx, reqCtx, in)
		assert.NoError(t, err, msg)
		assert.Equal(t, 2, len(out.Articles), msg+" - expected list with 2 articles")
	}

	{
		msg := "articles_list_sfl - with limit and offset"

		currUsername := username1

		reqCtx := web.RequestContext{
			Username: currUsername,
			Token:    &jwt.Token{},
		}

		in := rpc.ArticleCriteria{
			Author: util.PointerOf(username2),
			Limit:  util.PointerOf(1),
			Offset: util.PointerOf(1),
		}

		out, err := articlesListSfl(ctx, reqCtx, in)
		assert.NoError(t, err, msg)
		assert.Equal(t, 1, len(out.Articles), msg+" - expected list with one article")

		feedArticlce := out.Articles[0].Article.ToArticle()

		allArticles, err := testutil.ArticlesList(db, ctx, currUsername)
		assert.NoError(t, err, msg)
		assert.Equal(t, 2, len(allArticles), msg+" - retrieving all articles, expected list with 2 article")
		expectedArticle := allArticles[1]
		assert.Equal(t, expectedArticle, feedArticlce, msg+" - article from SFL must match second item in the list of all articles")
	}
}
