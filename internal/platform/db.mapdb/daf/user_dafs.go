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

type UserGetByNameDafT = func(userName string) (*model.User, error)

func (s UserDafs) getByName(userName string) (*model.User, error) {
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

func (s UserDafs) MakeGetByName() UserGetByNameDafT {
	return s.getByName
}
