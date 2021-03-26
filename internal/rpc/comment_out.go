package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type CommentOut struct {
	Comment model.Comment
}

type CommentsOut struct {
	Comments []CommentOut
}
