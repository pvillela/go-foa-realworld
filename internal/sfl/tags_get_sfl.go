/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// TagsGetSflT is the type of the stereotype instance for the service flow that
// retrieves all tags.
type TagsGetSflT = func(ctx context.Context, reqCtx web.RequestContext, _ types.Unit) (rpc.TagsOut, error)

// TagsGetSflC is the function that constructs a stereotype instance of type
// TagsGetSflT with hard-wired stereotype dependencies.

// TagsGetSflC0 is the function that constructs a stereotype instance of type
// TagsGetSflT without hard-wired stereotype dependencies.
func TagsGetSflC0(
	db dbpgx.Db,
	tagsGetAllDaf daf.TagsGetAllDafT,
) TagsGetSflT {
	return func(ctx context.Context, reqCtx web.RequestContext, _ types.Unit) (rpc.TagsOut, error) {
		return dbpgx.Db_WithTransaction(db, ctx, func(
			ctx context.Context,
			tx pgx.Tx,
		) (rpc.TagsOut, error) {
			tags, err := tagsGetAllDaf(ctx, tx)
			if err != nil {
				return rpc.TagsOut{}, err
			}

			tagsOut := rpc.TagsOut_FromModel(tags)

			return tagsOut, err
		})
	}
}
