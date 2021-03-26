package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type CommentOut struct {
	Comment *model.Comment
}

func (self CommentOut) FromModel(comment *model.Comment) CommentOut {
	self.Comment = comment
	return self
}

type CommentsOut struct {
	Comments []CommentOut
}

func (self CommentsOut) FromModel(comments []model.Comment) CommentsOut {
	outs := make([]CommentOut, len(comments))
	for i, comment := range comments {
		outs[i] = CommentOut{}.FromModel(&comment)
	}
	return CommentsOut{outs}
}
