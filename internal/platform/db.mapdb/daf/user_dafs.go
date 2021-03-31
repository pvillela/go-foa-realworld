package daf

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
)

type UserDafs struct {
	Store *sync.Map
}

func (s UserDafs) MakeGetByName() fs.UserGetByNameDafT {
	return func(userName string) (*model.User, error) {
		value, ok := s.Store.Load(userName)
		if !ok {
			return nil, fs.ErrUserNotFound
		}

		user, ok := value.(model.User)
		if !ok {
			return nil, errors.New("corrupted data, expected value of type User at key " + userName)
		}

		return &user, nil
	}
}

func (s UserDafs) MakeGetByEmail() fs.UserGetByEmailDafT {
	return func(email string) (*model.User, error) {
		var err error
		var foundUser model.User
		var userWasFound bool

		s.Store.Range(func(key, value interface{}) bool {
			user, ok := value.(model.User)
			if !ok {
				err = errors.New("data corrupton: expected model.User")
				return false
			}

			if user.Email == email {
				foundUser = user
				userWasFound = true
				return false // stop range
			}

			return true // keep iterating
		})

		var userRet *model.User
		if userWasFound {
			userRet = &foundUser
		}
		return userRet, nil
	}
}
