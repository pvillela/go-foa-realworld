package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/jwt"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type UserAuthenticateBf struct{}

type UserAuthenticateBfT = func(user *model.User, password string) bool

func (UserAuthenticateBf) Make() UserAuthenticateBfT {
	return func(user *model.User, password string) bool {
		if jwt.Hash(user.PasswordSalt, password) != user.PasswordHash {
			return false
		}
		return true
	}
}
