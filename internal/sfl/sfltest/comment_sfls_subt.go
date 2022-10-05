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
	"github.com/pvillela/go-foa-realworld/internal/rpc"
	"github.com/pvillela/go-foa-realworld/internal/sfl/boot"
	"github.com/pvillela/go-foa-realworld/internal/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

///////////////////
// Shared constants and data

type commentSourceT struct {
	username string
	slug     string
	body     string
}

var commentSources = []commentSourceT{
	{
		username: username1,
		slug:     util.Slug(authorsAndArticles[0].Article.Title),
		body:     "I thought it was a great article.",
	},
	{
		username: username1,
		slug:     util.Slug(authorsAndArticles[1].Article.Title),
		body:     "Not too bad.",
	},
	{
		username: username3,
		slug:     util.Slug(authorsAndArticles[1].Article.Title),
		body:     "What a terrible article.",
	},
}

///////////////////
// Tests

func commentAddSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	boot.CommentAddSflCfgAdapter = TestCfgAdapterOf(db)
	commentAddSfl := boot.CommentAddSflBoot(nil)

	{
		msg := "comment_add_sfl - valid comment"

		for _, c := range commentSources {
			authorname := c.username
			reqCtx := web.RequestContext{
				Username: authorname,
				Token:    &jwt.Token{},
			}

			in := rpc.CommentAddIn{
				Slug:    c.slug,
				Comment: rpc.CommentAddIn0{Body: util.PointerOf(c.body)},
			}

			commentOut, err := commentAddSfl(ctx, reqCtx, in)
			assert.NoError(t, err)

			returned := commentOut.Comment

			user, err := testutil.UserGetByName(db, ctx, authorname)
			assert.NoError(t, err)

			article, err := testutil.ArticleGetBySlug(db, ctx, user.Username, c.slug)

			assert.Equal(t, user.Id, returned.AuthorId, msg+" - AuthorId must match input")
			assert.Equal(t, article.Id, returned.ArticleId, msg+" - ArticleId must match input")
			assert.Equal(t, c.body, *returned.Body, msg+" - Body must match input")
		}
	}

	{
		msg := "comment_add_sfl - second comment by same author"

		for _, c := range commentSources {
			authorname := c.username
			reqCtx := web.RequestContext{
				Username: authorname,
				Token:    &jwt.Token{},
			}

			in := rpc.CommentAddIn{
				Slug:    c.slug,
				Comment: rpc.CommentAddIn0{Body: util.PointerOf(c.body + "2")},
			}

			commentOut, err := commentAddSfl(ctx, reqCtx, in)
			assert.NoError(t, err)

			returned := commentOut.Comment

			user, err := testutil.UserGetByName(db, ctx, authorname)
			assert.NoError(t, err)

			article, err := testutil.ArticleGetBySlug(db, ctx, user.Username, c.slug)

			assert.Equal(t, user.Id, returned.AuthorId, msg+" - AuthorId must match input")
			assert.Equal(t, article.Id, returned.ArticleId, msg+" - ArticleId must match input")
			assert.Equal(t, c.body+"2", *returned.Body, msg+" - Body must match input")
		}
	}

	{
		msg := "comment_add_sfl - inexistent article"

		authorname := username1
		reqCtx := web.RequestContext{
			Username: authorname,
			Token:    &jwt.Token{},
		}

		in := rpc.CommentAddIn{
			Slug:    "xxx",
			Comment: rpc.CommentAddIn0{Body: util.PointerOf("A comment.")},
		}

		_, err := commentAddSfl(ctx, reqCtx, in)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[article slug"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when username or email is not unique")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username or email is not unique")
	}
}

//func commentDeleteSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
//	boot.CommentDeleteSflCfgAdapter = TestCfgAdapterOf(db)
//	commentDeleteSfl := boot.CommentDeleteSflBoot(nil)
//
//	{
//		msg := "comment_delete_sfl - existing comment authored by current user"
//
//		aa := commentSources[2]
//		comment := aa.Comment
//		slug := util.Slug(comment.Title)
//		authorname := aa.Username
//
//		reqCtx := web.RequestContext{
//			Username: authorname,
//			Token:    &jwt.Token{},
//		}
//
//		_, err := commentDeleteSfl(ctx, reqCtx, slug)
//		assert.NoError(t, err, msg)
//
//		_, err = testutil.CommentsGetBySlug(db, ctx, authorname, slug)
//		returnedErrxKind := dbpgx.ClassifyError(err)
//		expectedErrxKind := dbpgx.DbErrRecordNotFound
//		expectedErrMsgPrefix := "DbErrRecordNotFound[comment slug"
//
//		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - retrieval of deleted comment must fail with appropriate error kind when username or email is not unique")
//		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - retrieval of deleted comment must fail with appropriate error message when username or email is not unique")
//	}
//
//	{
//		msg := "comment_delete_sfl - inexistenet comment"
//
//		aa := commentSources[2]
//		comment := aa.Comment
//		slug := util.Slug(comment.Title)
//		authorname := aa.Username
//
//		reqCtx := web.RequestContext{
//			Username: authorname,
//			Token:    &jwt.Token{},
//		}
//
//		_, err := commentDeleteSfl(ctx, reqCtx, slug)
//		returnedErrxKind := dbpgx.ClassifyError(err)
//		expectedErrxKind := dbpgx.DbErrRecordNotFound
//		expectedErrMsgPrefix := "DbErrRecordNotFound[comment slug"
//
//		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when username or email is not unique")
//		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username or email is not unique")
//	}
//
//	{
//		msg := "comment_delete_sfl - attempted by non-author"
//
//		aa := commentSources[0]
//		comment := aa.Comment
//		slug := util.Slug(comment.Title)
//		authorname := username1
//
//		reqCtx := web.RequestContext{
//			Username: authorname,
//			Token:    &jwt.Token{},
//		}
//
//		_, err := commentDeleteSfl(ctx, reqCtx, slug)
//		returnedErrxKind := dbpgx.ClassifyError(err)
//		expectedErrxKind := bf.ErrUnauthorizedUser
//		expectedErrMsgPrefix := "ErrUnauthorizedUser[user"
//
//		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when username or email is not unique")
//		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username or email is not unique")
//	}
//}

//func commentsGetSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
//	boot.CommentsGetSflCfgAdapter = TestCfgAdapterOf(db)
//	commentGetSfl := boot.CommentsGetSflBoot(nil)
//
//	{
//		msg := "comment_get_sfl - existing comment"
//
//		currUsername := username1
//		aa := commentSources[0]
//		comment := aa.Comment
//		slug := util.Slug(comment.Title)
//
//		reqCtx := web.RequestContext{
//			Username: currUsername,
//			Token:    &jwt.Token{},
//		}
//
//		out, err := commentGetSfl(ctx, reqCtx, slug)
//		assert.NoError(t, err, msg)
//
//		assert.True(t, out.Comment.Favorited, msg+" - Favorited attribute of output must be true because it already was before")
//		assert.Equal(t, comment.Title, out.Comment.Title, msg+" - Title attribute must not change")
//		assert.Equal(t, comment.Description, out.Comment.Description, msg+" - Description attribute must not change")
//		assert.Equal(t, comment.Body, out.Comment.Body, msg+" - Body attribute must not change")
//	}
//
//	{
//		msg := "comment_get_sfl - inexistent comment"
//
//		currUsername := username1
//		slug := "dkdkddkd"
//
//		reqCtx := web.RequestContext{
//			Username: currUsername,
//			Token:    &jwt.Token{},
//		}
//
//		_, err := commentGetSfl(ctx, reqCtx, slug)
//		returnedErrxKind := dbpgx.ClassifyError(err)
//		expectedErrxKind := dbpgx.DbErrRecordNotFound
//		expectedErrMsgPrefix := "DbErrRecordNotFound[comment slug"
//
//		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when favoriting an inexistent comment")
//		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when favoriting an inexistent comment")
//	}
//}
