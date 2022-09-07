/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import "github.com/pvillela/go-foa-realworld/internal/bf"

type ArticleUpdateIn struct {
	Article struct {
		Slug        string
		Title       *string // optional
		Description *string // optional
		Body        *string // optional
	}
}

func (in ArticleUpdateIn) Validate() error {
	if in.Article.Slug == "" {
		return bf.ErrValidationFailed.Make(nil,
			"article slug missing for Update operation")
	}
	return nil
}
