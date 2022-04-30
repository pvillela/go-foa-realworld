/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"

	"github.com/pvillela/go-foa-realworld/internal/arch"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentDeleteSflT is the type of the stereotype instance for the service flow that
// deletes a comment from an article.
type CommentDeleteSflT = func(ctx context.Context, in rpc.CommentDeleteIn) (arch.Unit, error)

// CommentDeleteSflC is the function that constructs a stereotype instance of type
// CommentDeleteSflT.
func CommentDeleteSflC(
	beginTxn func(context string) db.Txn,
	commentGetByIdDaf fs.CommentGetByIdDafT,
	commentDeleteDaf fs.CommentDeleteDafT,
	articleGetBySlugdDaf fs.ArticleGetBySlugDafT,
	articleUpdateDaf fs.ArticleUpdateDafT,
) CommentDeleteSflT {
	return func(ctx context.Context, in rpc.CommentDeleteIn) (arch.Unit, error) {
		username := web.ContextToRequestContext(ctx).Username

		txn := beginTxn("ArticleCreateSflS")
		defer txn.End()

		article, _, err := articleGetBySlugdDaf(in.Slug)
		if err != nil {
			return arch.Void, err
		}
		comment, _, err := commentGetByIdDaf(article.Id, in.Id)
		if err != nil {
			return arch.Void, err
		}
		if comment.Author.Username != username {
			return arch.Void, fs.ErrUnauthorizedUser.Make(nil, username)
		}

		if err := commentDeleteDaf(article.Id, in.Id, txn); err != nil {
			return arch.Void, err
		}

		article, rc, err := articleGetBySlugdDaf(in.Slug)
		if err != nil {
			return arch.Void, err
		}

		article = article.UpdateComments(comment, false)

		if _, err := articleUpdateDaf(article, rc, txn); err != nil {
			return arch.Void, err
		}

		return arch.Void, nil
	}
}
