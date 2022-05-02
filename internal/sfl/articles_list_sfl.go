/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticlesListSflT is the type of the stereotype instance for the service flow that
// retrieve recent articles based on a set of query parameters.
type ArticlesListSflT = func(ctx context.Context, in model.ArticleCriteria) (rpc.ArticlesOut, error)

// ArticlesListSflC is the function that constructs a stereotype instance of type
// ArticlesListSflT.
func ArticlesListSflC(
	db dbpgx.Db,
	articlesListDaf daf.ArticlesListDafT,
) ArticlesListSflT {
	return func(ctx context.Context, in model.ArticleCriteria) (rpc.ArticlesOut, error) {
		tx, err := db.BeginTx(ctx)
		if err != nil {
			return rpc.ArticlesOut{}, err
		}

		articles, err := articlesListDaf(ctx, tx, in)
		if err != nil {
			return rpc.ArticlesOut{}, err
		}

		articlesOut := rpc.ArticlesOut_FromModel(user, articles)
		return articlesOut, err
	}
}
