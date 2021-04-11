/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import "github.com/pvillela/go-foa-realworld/internal/model"

type CommentAddIn struct {
	Slug    string
	Comment commentAddIn0
}

type commentAddIn0 struct {
	Body *string
}

func (in CommentAddIn) ToComment(commentAuthor model.User) model.Comment {
	comment := model.Comment{
		Body:   in.Comment.Body,
		Author: commentAuthor,
	}
	return comment
}
