/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type CommentAddIn struct {
	Slug    string
	Comment commentAddIn0
}

type commentAddIn0 struct {
	Body *string // mandatory
}

func (in CommentAddIn) ToComment(articleId uint, commentAuthorId uint) model.Comment {
	return model.Comment{
		ArticleId: articleId,
		AuthorId:  commentAuthorId,
		Body:      in.Comment.Body,
	}
}

func (in CommentAddIn) Validate() error {
	if in.Slug == "" || in.Comment.Body == nil {
		return bf.ErrValidationFailed.Make(nil, "slug or body is missing")
	}
	return nil
}
