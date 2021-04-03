package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserAuthenticateSfl is the stereotype instance for the service flow that
// authenticates a user.
type UserAuthenticateSfl struct {
	UserGetByEmailDaf  fs.UserGetByEmailDafT
	UserAuthenticateBf fs.UserAuthenticateBfT
	UserGenTokenBf     fs.UserGenTokenBfT
}

// UserAuthenticateSflT is the function type instantiated by UserAuthenticateSfl.
type UserAuthenticateSflT = func(_username string, in rpc.UserAuthenticateIn) (rpc.UserOut, error)

func (s UserAuthenticateSfl) Make() UserAuthenticateSflT {
	return func(_username string, in rpc.UserAuthenticateIn) (rpc.UserOut, error) {
		var zero rpc.UserOut

		email := in.User.Email
		password := in.User.Password

		pwUser, err := s.UserGetByEmailDaf(email)
		if err != nil {
			return zero, err
		}
		user := pwUser.Entity()

		if !s.UserAuthenticateBf(*user, password) {
			return zero, fs.ErrAuthenticationFailed
		}

		token, err := s.UserGenTokenBf(*user)
		if err != nil {
			return zero, err
		}

		userOut := rpc.UserOut{}.FromModel(*user, token)
		return userOut, err
	}
}
