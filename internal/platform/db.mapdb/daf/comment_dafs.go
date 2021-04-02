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
	return func(id int) (fs.PwComment, error) {
		var zeroPwComment fs.PwComment

		value, ok := s.Store.Load(id)
		if !ok {
			return zeroPwComment, fs.ErrCommentNotFound
		}

		pwComment, ok := value.(fs.PwComment)
		if !ok {
			return zeroPwComment, errors.New("not an article stored at key")
		}

		return pwComment, nil
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
	return func(comment model.Comment) (fs.PwComment, error) {
		comment.ID = s.getNextId()
		comment.CreatedAt = time.Now()
		comment.UpdatedAt = time.Now()
		pwComment := fs.PwComment{nil, comment}
		s.Store.Store(comment.ID, pwComment)
		return pwComment, nil
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
