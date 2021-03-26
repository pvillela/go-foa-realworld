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
