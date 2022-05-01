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

// ArticlesListSflT is the type of the stereotype instance for the service flow that
// retrieve recent articles based on a set of query parameters.
type ArticlesListSflT = func(ctx context.Context, in rpc.ArticlesListIn) (rpc.ArticlesOut, error)

// ArticlesListSflC is the function that constructs a stereotype instance of type
// ArticlesListSflT.
func ArticlesListSflC(
	userGetByNameDaf newdaf.UserGetByNameDafT,
	articleGetRecentFilteredDaf fs.ArticleGetRecentFilteredDafT,
) ArticlesListSflT {
	return func(ctx context.Context, in rpc.ArticlesListIn) (rpc.ArticlesOut, error) {
		username := web.ContextToRequestContext(ctx).Username

		var zero rpc.ArticlesOut
		var user model.User
		var err error

		if username != "" {
			user, _, err = userGetByNameDaf(username)
			if err != nil {
				return zero, err
			}
		}

		articles, err := articleGetRecentFilteredDaf(in)
		if err != nil {
			return zero, err
		}

		articlesOut := rpc.ArticlesOut_FromModel(user, articles)
		return articlesOut, err
	}
}
