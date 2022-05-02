/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleGetSflT is the type of the stereotype instance for the service flow that
// retrieves an article.
type ArticleGetSflT = func(ctx context.Context, slug string) (rpc.ArticleOut, error)

// ArticleGetSflC is the function that constructs a stereotype instance of type
// ArticleGetSflT.
func ArticleGetSflC(
	ctxDb db.CtxDb,
	userGetByNameDaf daf.UserGetByNameDafT,
	articleGetBySlugDaf daf.ArticleGetBySlugDafT,
) ArticleGetSflT {
	return func(ctx context.Context, slug string) (rpc.ArticleOut, error) {
		ctx, err := ctxDb.BeginTx(ctx)
		if err != nil {
			return rpc.ArticleOut{}, err
		}
		defer ctxDb.DeferredRollback(ctx)

		article, _, err := articleGetBySlugDaf(slug)
		if err != nil {
			return rpc.ArticleOut{}, err
		}

		username := web.ContextToRequestContext(ctx).Username
		var user model.User
		if username != "" {
			user, _, err = userGetByNameDaf(ctx, username)
			if err != nil {
				return rpc.ArticleOut{}, err
			}
		}

		articleOut := rpc.ArticleOut_FromModel(article, followsAuthor, likesArticle)

		return articleOut, err
	}
}
