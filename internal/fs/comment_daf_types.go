package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type MdbComment struct {
	db.RecCtx
	Entity model.Comment
}

type CommentGetByIdDafT = func(id int) (MdbComment, error)

type CommentCreateDafT = func(comment model.Comment) (MdbComment, error)

type CommentDeleteDafT = func(id int) error
