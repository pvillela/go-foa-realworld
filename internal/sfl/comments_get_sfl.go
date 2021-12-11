/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentsGetSfl is the stereotype instance for the service flow that
// retrieves the comments of an article.
type CommentsGetSfl struct {
	ArticleGetBySlugDaf fs.ArticleGetBySlugDafT
}

// CommentsGetSflT is the function type instantiated by CommentsGetSfl.
type CommentsGetSflT = func(username string, slug string) (rpc.CommentsOut, error)

func (s CommentsGetSfl) Make() CommentsGetSflT {
	return func(username string, slug string) (rpc.CommentsOut, error) {
		var zero rpc.CommentsOut

		article, _, err := s.ArticleGetBySlugDaf(slug)
		if err != nil {
			return zero, err
		}

		if article.Comments == nil {
			article.Comments = []model.Comment{}
		}

		commentsOut := rpc.CommentsOut_FromModel(article.Comments)

		return commentsOut, nil
	}
}
