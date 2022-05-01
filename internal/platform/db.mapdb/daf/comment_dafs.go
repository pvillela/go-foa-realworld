/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/newdaf"
	"strconv"
	"strings"

	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

func commentKey(comment model.Comment) string {
	return commentKey0(comment.ArticleId, comment.Id)
}

func commentKey0(articleUuid util.Uuid, id int) string {
	return string(articleUuid) + "-" + strconv.Itoa(id)
}

func pwCommentFromDb(value interface{}) newdaf.PwComment {
	pw, ok := value.(newdaf.PwComment)
	if !ok {
		panic(fmt.Sprintln("database corrupted, value", pw, "does not wrap comment"))
	}
	return pw
}

// CommentGetByIdDafC is the function that constructs a stereotype instance of type
// fs.CommentGetByIdDafT.
func CommentGetByIdDafC(commentDb mapdb.MapDb) newdaf.CommentGetByIdDafT {
	return func(articleUuid util.Uuid, id int) (model.Comment, newdaf.RecCtxComment, error) {
		value, err := commentDb.Read(commentKey0(articleUuid, id))
		if err != nil {
			return model.Comment{}, newdaf.RecCtxComment{}, fs.ErrCommentNotFound.Make(err, articleUuid, id)
		}
		pw := pwCommentFromDb(value)
		return pw.Entity, pw.RecCtx, nil
	}
}

func nextId(commentDb mapdb.MapDb, articleUuid util.Uuid) int {
	nextId := 1
	commentDb.Range(func(key, value interface{}) bool {
		if strings.HasPrefix(key.(string), string(articleUuid)) {
			nextId++
		}
		return false
	})
	return nextId
}

// CommentCreateDafC is the function that constructs a stereotype instance of type
// fs.CommentCreateDafT.
func CommentCreateDafC(commentDb mapdb.MapDb) newdaf.CommentCreateDafT {
	return func(comment model.Comment, txn db.Txn) (model.Comment, newdaf.RecCtxComment, error) {
		comment.Id = nextId(commentDb, comment.ArticleId)
		pw := newdaf.PwComment{newdaf.RecCtxComment{}, comment}
		if err := commentDb.Create(commentKey(comment), pw, txn); err != nil {
			return model.Comment{}, newdaf.RecCtxComment{}, err // can only be an invalid txn token due to first line above
		}

		return comment, newdaf.RecCtxComment{}, nil
	}
}

// CommentDeleteDafC is the function that constructs a stereotype instance of type
// fs.CommentDeleteDafT.
func CommentDeleteDafC(commentDb mapdb.MapDb) newdaf.CommentDeleteDafT {
	return func(articleUuid util.Uuid, id int, txn db.Txn) error {
		err := commentDb.Delete(commentKey0(articleUuid, id), txn)
		if errx.KindOf(err) == mapdb.ErrRecordNotFound {
			return fs.ErrCommentNotFound.Make(err, articleUuid, id)
		}
		return err
	}
}
