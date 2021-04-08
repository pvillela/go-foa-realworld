package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// PwUser is a wrapper of the model.User entity
// containing context information required for ersistence purposes.
type PwUser struct {
	db.RecCtx
	Entity model.User
}

type UserGetByNameDafT = func(userName string) (model.User, db.RecCtx, error)

type UserGetByEmailDafT = func(email string) (model.User, db.RecCtx, error)

type UserUpdateDafT = func(user model.User, recCtx db.RecCtx) (model.User, db.RecCtx, error)
