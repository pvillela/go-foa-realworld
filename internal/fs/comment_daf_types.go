package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// PwComment is a wrapper of the model.User entity
// containing context information required for ersistence purposes.
type PwComment struct {
	db.RecCtx
	Entity model.Comment
}

type CommentGetByIdDafT = func(id int) (PwComment, error)

type CommentCreateDafT = func(comment model.Comment) (PwComment, error)

type CommentDeleteDafT = func(id int) error
