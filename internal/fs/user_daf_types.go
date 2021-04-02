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

type UserGetByNameDafT = func(userName string) (PwUser, error)

type UserGetByEmailDafT = func(email string) (PwUser, error)

type UserUpdateDafT = func(pwUser PwUser) (PwUser, error)
