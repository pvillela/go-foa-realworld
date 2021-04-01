package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type MdbUser struct {
	db.RecCtx
	Entity model.User
}

type UserGetByNameDafT = func(userName string) (MdbUser, error)

type UserGetByEmailDafT = func(email string) (MdbUser, error)

type UserUpdateDafT = func(mdbUser MdbUser) (MdbUser, error)
