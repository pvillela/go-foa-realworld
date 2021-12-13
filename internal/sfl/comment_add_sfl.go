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

// CommentAddSflT is the type of the stereotype instance for the service flow that
// adds a comment to an article.
type CommentAddSflT = func(username string, in rpc.CommentAddIn) (rpc.CommentOut, error)

// CommentAddSflC is the function that constructs a stereotype instance of type
// CommentAddSflT.
func CommentAddSflC(
	beginTxn func(context string) db.Txn,
	userGetByNameDaf fs.UserGetByNameDafT,
	articleGetBySlugDaf fs.ArticleGetBySlugDafT,
	commentCreateDaf fs.CommentCreateDafT,
	articleUpdateDaf fs.ArticleUpdateDafT,
) CommentAddSflT {
	return func(username string, in rpc.CommentAddIn) (rpc.CommentOut, error) {
		txn := beginTxn("ArticleCreateSflS")
		defer txn.End()

		var zero rpc.CommentOut
		var err error

		commentAuthor, _, err := userGetByNameDaf(username)
		if err != nil {
			return zero, err
		}

		article, rc, err := articleGetBySlugDaf(in.Slug)
		if err != nil {
			return zero, err
		}

		rawComment := in.ToComment(article.Uuid, commentAuthor)

		insertedComment, _, err := commentCreateDaf(rawComment, txn)
		if err != nil {
			return zero, err
		}

		article.Comments = append(article.Comments, insertedComment)

		if _, err := articleUpdateDaf(article, rc, txn); err != nil {
			return zero, err
		}

		commentOut := rpc.CommentOut_FromModel(insertedComment)
		return commentOut, err
	}
}
