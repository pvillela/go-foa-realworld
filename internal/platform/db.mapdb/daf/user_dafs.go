package daf

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
	"time"
)

type UserDafs struct {
	Store *sync.Map
}

func (s UserDafs) getByName(username string) (*model.User, error) {
	value, ok := s.Store.Load(username)
	if !ok {
		return nil, nil, fs.ErrUserNotFound
	}

	user, ok := value.(model.User)
	if !ok {
		return nil, nil, errors.New("corrupted data, expected value of type Entity at key " + username)
	}

	return &user, nil, nil
}

func (s UserDafs) MakeGetByName() fs.UserGetByNameDafT {
	return s.getByName
}

func (s UserDafs) MakeGetByEmail() fs.UserGetByEmailDafT {
	return func(email string) (*model.User, error) {
		var err error
		var foundUser model.User
		var userWasFound bool

		s.Store.Range(func(key, value interface{}) bool {
			user, ok := value.(model.User)
			if !ok {
				err = errors.New("data corrupton: expected model.Entity")
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

func (s UserDafs) MakeUpdate() fs.UserUpdateDafT {
	return func(user model.User) error {
		if user, _ := s.getByName(user.Name); user == nil {
			return fs.ErrUserNotFound
		}

		user.UpdatedAt = time.Now()
		s.Store.Store(user.Name, user)

		return nil
	}
}
