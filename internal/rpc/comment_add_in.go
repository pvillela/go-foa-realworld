package rpc

import "github.com/pvillela/go-foa-realworld/internal/model"

type CommentAddIn struct {
	Slug    string
	Comment CommentAddIn0
}

type CommentAddIn0 struct {
	Body string
}

func (in CommentAddIn) ToComment(commentAuthor *model.User) model.Comment {
	comment := model.Comment{
		Body:   in.Comment.Body,
		Author: *commentAuthor,
	}
	return comment
}
