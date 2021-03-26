package daf

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
	"time"
)

type CommentDafs struct {
	Store *sync.Map
}

func (s CommentDafs) MakeGetById() fs.CommentGetByIdDafT {
	return func(id int) (*model.Comment, error) {
		value, ok := s.Store.Load(id)
		if !ok {
			return nil, fs.ErrCommentNotFound
		}

		comment, ok := value.(model.Comment)
		if !ok {
			return nil, errors.New("not an article stored at key")
		}

		return &comment, nil
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
	return func(comment model.Comment) (*model.Comment, error) {
		comment.ID = s.getNextId()
		comment.CreatedAt = time.Now()
		comment.UpdatedAt = time.Now()
		s.Store.Store(comment.ID, comment)
		return &comment, nil
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
