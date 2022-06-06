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

	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentAddSflT is the type of the stereotype instance for the service flow that
// adds a comment to an article.
type CommentAddSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	in rpc.CommentAddIn,
) (rpc.CommentOut, error)

// CommentAddSflC is the function that constructs a stereotype instance of type
// CommentAddSflT with hard-wired stereotype dependencies.
func CommentAddSflC(
	db dbpgx.Db,
) CommentAddSflT {
	articleAndUserGetFl := fl.ArticleAndUserGetFlI
	commentCreateDaf := daf.CommentCreateDafI
	return CommentAddSflC0(
		db,
		articleAndUserGetFl,
		commentCreateDaf,
	)
}

// CommentAddSflC0 is the function that constructs a stereotype instance of type
// CommentAddSflT without hard-wired stereotype dependencies.
func CommentAddSflC0(
	db dbpgx.Db,
	articleAndUserGetFl fl.ArticleAndUserGetFlT,
	commentCreateDaf daf.CommentCreateDafT,
) CommentAddSflT {
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		in rpc.CommentAddIn,
	) (rpc.CommentOut, error) {
		err := in.Validate()
		if err != nil {
			return rpc.CommentOut{}, err
		}

		username := reqCtx.Username

		articlePlus, user, err := articleAndUserGetFl(ctx, tx, in.Slug, username)
		if err != nil {
			return rpc.CommentOut{}, err
		}

		comment := in.ToComment(articlePlus.Id, user.Id)

		err = commentCreateDaf(ctx, tx, &comment)
		if err != nil {
			return rpc.CommentOut{}, err
		}

		commentOut := rpc.CommentOut_FromModel(comment)
		return commentOut, nil
	})
}
