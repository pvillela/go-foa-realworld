package daf

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
	"time"
)

type CommentDafs struct {
	Store *sync.Map
}

func (s CommentDafs) MakeGetById() fs.CommentGetByIdDafT {
	return func(id int) (model.Comment, db.RecCtx, error) {
		value, ok := s.Store.Load(id)
		if !ok {
			return model.Comment{}, nil, fs.ErrCommentNotFound
		}

		pw, ok := value.(fs.PwComment)
		if !ok {
			panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap comment"))
		}

		return pw.Entity, pw.RecCtx, nil
	}
}

func (s CommentDafs) getNextId() int {
	nextId := 0
	s.Store.Range(func(key, value interface{}) bool {
		nextId++
		return true
	})
	return nextId
}

func (s CommentDafs) MakeCreate() fs.CommentCreateDafT {
	return func(comment model.Comment) (model.Comment, db.RecCtx, error) {
		comment.ID = s.getNextId()
		comment.CreatedAt = time.Now()
		comment.UpdatedAt = time.Now()
		pw := fs.PwComment{nil, comment}
		s.Store.Store(comment.ID, pw)
		return pw.Entity, nil, nil
	}
}

func (s CommentDafs) MakeDelete() fs.CommentDeleteDafT {
	return func(id int) error {
		if _, present := s.Store.LoadAndDelete(id); !present {
			return fs.ErrCommentNotFound
		}
		return nil
	}
}
