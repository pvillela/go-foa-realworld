/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentDeleteSfl is the stereotype instance for the service flow that
// deletes a comment from an article.
type CommentDeleteSfl struct {
	CommentGetByIdDaf    fs.CommentGetByIdDafT
	CommentDeleteDaf     fs.CommentDeleteDafT
	ArticleGetBySlugdDaf fs.ArticleGetBySlugDafT
	ArticleUpdateDaf     fs.ArticleUpdateDafT
}

// CommentDeleteSflT is the function type instantiated by CommentDeleteSfl.
type CommentDeleteSflT = func(username string, in rpc.CommentDeleteIn) error

func (s CommentDeleteSfl) Make() CommentDeleteSflT {
	return func(username string, in rpc.CommentDeleteIn) error {
		comment, _, err := s.CommentGetByIdDaf(in.Id)
		if err != nil {
			return err
		}
		if comment.Author.Name != username {
			return fs.ErrUnauthorizedUser
		}

		if err := s.CommentDeleteDaf(in.Id); err != nil {
			return err
		}

		article, rc, err := s.ArticleGetBySlugdDaf(in.Slug)
		if err != nil {
			return err
		}

		article = article.UpdateComments(comment, false)

		if _, _, err := s.ArticleUpdateDaf(article, rc); err != nil {
			return err
		}

		return nil
	}
}
