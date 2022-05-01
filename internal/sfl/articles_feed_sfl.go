/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/newdaf"

	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticlesFeedSflT is the type of the stereotype instance for the service flow that
// queries for all articles authored by other users followed by
// the current user, with optional limit and offset pagination parameters.
type ArticlesFeedSflT = func(ctx context.Context, in rpc.ArticlesFeedIn) (rpc.ArticlesOut, error)

// ArticlesFeedSflC is the function that constructs a stereotype instance of type
// ArticlesFeedSflT.
func ArticlesFeedSflC(
	userGetByNameDaf newdaf.UserGetByNameDafT,
	articleGetRecentForAuthorsDaf fs.ArticleGetRecentForAuthorsDafT,
) ArticlesFeedSflT {
	return func(ctx context.Context, in rpc.ArticlesFeedIn) (rpc.ArticlesOut, error) {
		username := web.ContextToRequestContext(ctx).Username

		var zero rpc.ArticlesOut
		var user model.User
		var err error

		if username == "" {
			return zero, fs.ErrNotAuthenticated.Make(nil)
		}

		user, _, err = userGetByNameDaf(username)
		if err != nil {
			return zero, err
		}

		articles, err := articleGetRecentForAuthorsDaf(user.Followees, in.Limit, in.Offset)
		if err != nil {
			return zero, err
		}

		articlesOut := rpc.ArticlesOut_FromModel(user, articles)
		return articlesOut, err
	}
}
