package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/jwt"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type UserGenTokenBf struct{}

type UserGenTokenBfT = func(user model.User) (string, error)

func (UserGenTokenBf) Make() UserGenTokenBfT {
	return func(user model.User) (string, error) {
		return jwt.UserGenToken(user)
	}
}
