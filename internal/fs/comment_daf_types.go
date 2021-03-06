/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// PwComment is a wrapper of the model.User entity
// containing context information required for persistence purposes.
type PwComment struct {
	db.RecCtx
	Entity model.Comment
}

type CommentGetByIdDafT = func(articleUuid util.Uuid, id int) (model.Comment, db.RecCtx, error)

type CommentCreateDafT = func(comment model.Comment, txn db.Txn) (model.Comment, db.RecCtx, error)

type CommentDeleteDafT = func(articleUuid util.Uuid, id int, txn db.Txn) error
