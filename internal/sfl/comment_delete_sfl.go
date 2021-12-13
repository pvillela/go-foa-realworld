/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentDeleteSflT is the type of the stereotype instance for the service flow that
// deletes a comment from an article.
type CommentDeleteSflT = func(username string, in rpc.CommentDeleteIn) error

// CommentDeleteSflC is the function that constructs a stereotype instance of type
// CommentDeleteSflT.
func CommentDeleteSflC(
	beginTxn func(context string) db.Txn,
	commentGetByIdDaf fs.CommentGetByIdDafT,
	commentDeleteDaf fs.CommentDeleteDafT,
	articleGetBySlugdDaf fs.ArticleGetBySlugDafT,
	articleUpdateDaf fs.ArticleUpdateDafT,
) CommentDeleteSflT {
	return func(username string, in rpc.CommentDeleteIn) error {
		txn := beginTxn("ArticleCreateSflS")
		defer txn.End()

		article, _, err := articleGetBySlugdDaf(in.Slug)
		if err != nil {
			return err
		}
		comment, _, err := commentGetByIdDaf(article.Uuid, in.Id)
		if err != nil {
			return err
		}
		if comment.Author.Name != username {
			return fs.ErrUnauthorizedUser.Make(nil, username)
		}

		if err := commentDeleteDaf(article.Uuid, in.Id, txn); err != nil {
			return err
		}

		article, rc, err := articleGetBySlugdDaf(in.Slug)
		if err != nil {
			return err
		}

		article = article.UpdateComments(comment, false)

		if _, err := articleUpdateDaf(article, rc, txn); err != nil {
			return err
		}

		return nil
	}
}
