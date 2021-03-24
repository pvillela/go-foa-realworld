package ft

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type UserGetByNameDafT = func(userName string) (*model.User, error)
