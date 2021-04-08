package daf

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
	"time"
)

type UserDafs struct {
	Store *sync.Map
}

func (s UserDafs) getByName(username string) (model.User, db.RecCtx, error) {
	value, ok := s.Store.Load(username)
	if !ok {
		return model.User{}, nil, fs.ErrUserNotFound
	}

	pw, ok := value.(fs.PwUser)
	if !ok {
		panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap pw"))
	}

	return pw.Entity, pw.RecCtx, nil
}

func (s UserDafs) MakeGetByName() fs.UserGetByNameDafT {
	return s.getByName
}

func (s UserDafs) MakeGetByEmail() fs.UserGetByEmailDafT {
	return func(email string) (model.User, db.RecCtx, error) {
		var foundPw fs.PwUser
		var pwWasFound bool

		s.Store.Range(func(key, value interface{}) bool {
			pw, ok := value.(fs.PwUser)
			if !ok {
				panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap pw"))
			}

			if pw.Entity.Email == email {
				foundPw = pw
				pwWasFound = true
				return false // stop range
			}

			return true // keep iterating
		})

		if !pwWasFound {
			return model.User{}, nil, fs.ErrUserNotFound
		}
		return foundPw.Entity, foundPw.RecCtx, nil
	}
}

func (s UserDafs) MakeUpdate() fs.UserUpdateDafT {
	return func(user model.User, recCtx db.RecCtx) (model.User, db.RecCtx, error) {
		if _, _, err := s.getByName(user.Name); err != nil {
			return model.User{}, nil, err
		}

		user.UpdatedAt = time.Now()
		pw := fs.PwUser{nil, user}
		s.Store.Store(user.Name, pw)

		return pw.Entity, pw.RecCtx, nil
	}
}
