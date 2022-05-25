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
	"github.com/pvillela/go-foa-realworld/internal/fl"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

	"github.com/pvillela/go-foa-realworld/internal/arch"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentDeleteSflT is the type of the stereotype instance for the service flow that
// deletes a comment from an article.
type CommentDeleteSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	in rpc.CommentDeleteIn,
) (arch.Unit, error)

// CommentDeleteSflC is the function that constructs a stereotype instance of type
// CommentDeleteSflT with hard-wired stereotype dependencies.
func CommentDeleteSflC(
	db dbpgx.Db,
) CommentDeleteSflT {
	articleAndUserGetFl := fl.ArticleAndUserGetFlI
	commentDeleteDaf := daf.CommentDeleteDafI
	return CommentDeleteSflC0(
		db,
		articleAndUserGetFl,
		commentDeleteDaf,
	)
}

// CommentDeleteSflC0 is the function that constructs a stereotype instance of type
// CommentDeleteSflT without hard-wired stereotype dependencies.
func CommentDeleteSflC0(
	db dbpgx.Db,
	articleAndUserGetFl fl.ArticleAndUserGetFlT,
	commentDeleteDaf daf.CommentDeleteDafT,
) CommentDeleteSflT {
	return func(
		ctx context.Context,
		reqCtx web.RequestContext,
		in rpc.CommentDeleteIn,
	) (arch.Unit, error) {
		return dbpgx.Db_WithTransaction(db, ctx, func(
			ctx context.Context,
			tx pgx.Tx,
		) (arch.Unit, error) {
			err := in.Validate()
			if err != nil {
				return arch.Void, err
			}

			username := reqCtx.Username

			articlePlus, user, err := articleAndUserGetFl(ctx, tx, in.Slug, username)
			if err != nil {
				return arch.Void, err
			}

			err = commentDeleteDaf(ctx, tx, uint(in.Id), articlePlus.Id, user.Id)
			if err != nil {
				return arch.Void, err
			}

			return arch.Void, nil
		})
	}
}
