/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentAddSfl is the stereotype instance for the service flow that
// adds a comment to an article.
type CommentAddSfl struct {
	BeginTxn            func(context string) db.Txn
	UserGetByNameDaf    fs.UserGetByNameDafT
	ArticleGetBySlugDaf fs.ArticleGetBySlugDafT
	CommentCreateDaf    fs.CommentCreateDafT
	ArticleUpdateDaf    fs.ArticleUpdateDafT
}

// CommentAddSflT is the function type instantiated by CommentAddSfl.
type CommentAddSflT = func(username string, in rpc.CommentAddIn) (rpc.CommentOut, error)

func (s CommentAddSfl) Make() CommentAddSflT {
	return func(username string, in rpc.CommentAddIn) (rpc.CommentOut, error) {
		txn := s.BeginTxn("ArticleCreateSfl")
		defer txn.End()

		var zero rpc.CommentOut
		var err error

		commentAuthor, _, err := s.UserGetByNameDaf(username)
		if err != nil {
			return zero, err
		}

		article, rc, err := s.ArticleGetBySlugDaf(in.Slug)
		if err != nil {
			return zero, err
		}

		rawComment := in.ToComment(article.Uuid, commentAuthor)

		insertedComment, _, err := s.CommentCreateDaf(rawComment, txn)
		if err != nil {
			return zero, err
		}

		article.Comments = append(article.Comments, insertedComment)

		if _, err := s.ArticleUpdateDaf(article, rc, txn); err != nil {
			return zero, err
		}

		commentOut := rpc.CommentOut{}.FromModel(insertedComment)
		return commentOut, err
	}
}
