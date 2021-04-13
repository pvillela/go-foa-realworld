/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package daf

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"sync"
)

type CommentDafs struct {
	CommentDb mapdb.MapDb
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
		pw := fs.PwComment{nil, comment}
		_, loaded := s.Store.LoadOrStore(comment.ID, pw)
		if loaded {
			return model.Comment{}, nil, fs.ErrDuplicateArticleSlug
		}

		return pw.Entity, pw.RecCtx, nil
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
