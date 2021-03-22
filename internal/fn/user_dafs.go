package fn

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
)

type UserDafS struct {
	Store *sync.Map
}

func (s UserDafS) GetByName(userName string) (*model.User, error) {
	value, ok := s.Store.Load(userName)
	if !ok {
		return nil, ErrUserNotFound
	}

	user, ok := value.(model.User)
	if !ok {
		return nil, errors.New("corrupted data, expected value of type User at key " + userName)
	}

	return &user, nil
}
