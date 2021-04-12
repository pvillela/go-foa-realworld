/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type CommentAddIn struct {
	Slug    string
	Comment commentAddIn0
}

type commentAddIn0 struct {
	Body *string // mandatory
}

// ToComment
// TODO: move this directly to the SFL because we have to read the article based on the slug
//   to get the uuid.
func (in CommentAddIn) ToComment(articleUuid util.Uuid, commentAuthor model.User) model.Comment {
	return model.Comment{}.Create(articleUuid, in.Comment.Body, commentAuthor)
}
