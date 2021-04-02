package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// PwComment is a wrapper of the model.User entity
// containing context information required for ersistence purposes.
type PwComment interface {
	Entity() *model.Comment
	SetEntity(*model.Comment)
	Copy(*model.Comment) PwComment
}

type CommentGetByIdDafT = func(id int) (PwComment, error)

type CommentCreateDafT = func(comment model.Comment) (PwComment, error)

type CommentDeleteDafT = func(id int) error
