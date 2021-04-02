package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// PwUser is a wrapper of the model.User entity
// containing context information required for ersistence purposes.
type PwUser interface {
	Entity() *model.User
	SetEntity(*model.User)
	Copy(*model.User) PwUser
}

type UserGetByNameDafT = func(userName string) (PwUser, error)

type UserGetByEmailDafT = func(email string) (PwUser, error)

type UserUpdateDafT = func(pwUser PwUser) (PwUser, error)
