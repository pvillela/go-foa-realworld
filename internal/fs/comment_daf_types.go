/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fs

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// RecCtxArticle is a type alias
type RecCtxComment = db.RecCtx[model.Comment]

// PwArticle is a type alias
type PwComment = db.Pw[model.Comment, RecCtxComment]

// CommentGetByIdDafT is the type of the stereotype instance for the DAF that
// retrieves a comment by Id.
type CommentGetByIdDafT = func(articleUuid util.Uuid, id int) (model.Comment, RecCtxComment, error)

// CommentCreateDafT is the type of the stereotype instance for the DAF that
// creates a comment.
type CommentCreateDafT = func(ctx context.Context, comment model.Comment) (model.Comment, RecCtxComment, error)

// CommentDeleteDafT is the type of the stereotype instance for the DAF that
// deletes a comment.
type CommentDeleteDafT = func(ctx context.Context, articleUuid util.Uuid, id int) error
