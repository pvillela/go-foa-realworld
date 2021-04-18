/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package daf

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"strconv"
	"strings"
)

type CommentDafs struct {
	CommentDb mapdb.MapDb
}

func commentKey(comment model.Comment) string {
	return commentKey0(comment.ArticleUuid, comment.ID)
}

func commentKey0(articleUuid util.Uuid, id int) string {
	return string(articleUuid) + "-" + strconv.Itoa(id)
}

func pwCommentFromDb(value interface{}) fs.PwComment {
	pw, ok := value.(fs.PwComment)
	if !ok {
		panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap comment"))
	}
	return pw
}

func (s CommentDafs) MakeGetByKey() fs.CommentGetByIdDafT {
	return func(articleUuid util.Uuid, id int) (model.Comment, db.RecCtx, error) {
		value, err := s.CommentDb.Read(commentKey0(articleUuid, id))
		if err != nil {
			return model.Comment{}, nil, fs.ErrCommentNotFound.Make(err, articleUuid, id)
		}
		pw := pwCommentFromDb(value)
		return pw.Entity, pw.RecCtx, nil
	}
}

func (s CommentDafs) nextId(articleUuid util.Uuid) int {
	nextId := 1
	s.CommentDb.Range(func(key, value interface{}) bool {
		if strings.HasPrefix(key.(string), string(articleUuid)) {
			nextId++
		}
		return false
	})
	return nextId
}

func (s CommentDafs) MakeCreate() fs.CommentCreateDafT {
	return func(comment model.Comment, txn db.Txn) (model.Comment, db.RecCtx, error) {
		comment.ID = s.nextId(comment.ArticleUuid)
		pw := fs.PwComment{nil, comment}
		if err := s.CommentDb.Create(commentKey(comment), pw, txn); err != nil {
			return model.Comment{}, nil, err // can only be an invalid txn token due to first line above
		}

		return comment, nil, nil
	}
}

func (s CommentDafs) MakeDelete() fs.CommentDeleteDafT {
	return func(articleUuid util.Uuid, id int, txn db.Txn) error {
		err := s.CommentDb.Delete(commentKey0(articleUuid, id), txn)
		if errx.KindOf(err) == mapdb.ErrRecordNotFound {
			return fs.ErrCommentNotFound.Make(err, articleUuid, id)
		}
		return err
	}
}
