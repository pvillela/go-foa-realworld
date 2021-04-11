/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type CommentOut struct {
	Comment model.Comment
}

func (s CommentOut) FromModel(comment model.Comment) CommentOut {
	s.Comment = comment
	return s
}

type CommentsOut struct {
	Comments []CommentOut
}

func (CommentsOut) FromModel(comments []model.Comment) CommentsOut {
	outs := make([]CommentOut, len(comments))
	for i, comment := range comments {
		outs[i] = CommentOut{}.FromModel(comment)
	}
	return CommentsOut{outs}
}
