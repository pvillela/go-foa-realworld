package fs

import "github.com/pvillela/go-foa-realworld/internal/model"

type CommentGetByIdDafT = func(id int) (*model.Comment, error)

type CommentCreateDafT = func(comment model.Comment) (*model.Comment, error)

type CommentDeleteDafT = func(id int) error
