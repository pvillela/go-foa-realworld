/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import "github.com/pvillela/go-foa-realworld/experimental/bf"

type CommentDeleteIn struct {
	Slug string
	Id   int
}

func (in CommentDeleteIn) Validate() error {
	if in.Slug == "" || in.Id == 0 {
		return bf.ErrValidationFailed.Make(nil, "article slug or comment id is missing")
	}
	return nil
}
