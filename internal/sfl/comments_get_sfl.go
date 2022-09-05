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
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/pvillela/go-foa-realworld/rpc"
)

// CommentsGetSflC0 is the type of the stereotype instance for the service flow that
// retrieves the comments of an article.
type CommentsGetSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	slug string,
) (rpc.CommentsOut, error)

// CommentsGetSflC is the function that constructs a stereotype instance of type
// CommentsGetSflT with hard-wired stereotype dependencies.
func CommentsGetSflC(
	cfgSrc DefaultSflCfgSrc,
) CommentsGetSflT {
	return CommentsGetSflC0(
		cfgSrc,
		daf.CommentsGetBySlugDaf,
	)
}

// CommentsGetSflC0 is the function that constructs a stereotype instance of type
// CommentsGetSflT without hard-wired stereotype dependencies.
func CommentsGetSflC0(
	cfgSrc DefaultSflCfgSrc,
	commentsGetBySlugDaf daf.CommentsGetBySlugDafT,
) CommentsGetSflT {
	db := cfgSrc()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		slug string,
	) (rpc.CommentsOut, error) {
		var zero rpc.CommentsOut

		comments, err := commentsGetBySlugDaf(ctx, tx, slug)
		if err != nil {
			return zero, err
		}

		commentsOut := rpc.CommentsOut_FromModel(comments)

		return commentsOut, nil
	})
}
