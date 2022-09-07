/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import (
	"github.com/pvillela/go-foa-realworld/experimental/bf"
	"github.com/pvillela/go-foa-realworld/experimental/model"
)

type ArticleCreateIn struct {
	Article ArticleCreateIn0
}

type ArticleCreateIn0 struct {
	Title       string
	Description string
	Body        *string   // mandatory
	TagList     *[]string // optional
}

func (in ArticleCreateIn) Validate() error {
	if in.Article.Title == "" || in.Article.Description == "" || in.Article.Body == nil {
		return bf.ErrValidationFailed.Make(nil,
			"article has missing fields for Create operation")
	}
	return nil
}

func (in ArticleCreateIn) ToArticle(author model.User) (model.Article, error) {
	tagList := in.Article.TagList
	if tagList == nil {
		tagList = new([]string)
	}

	return model.Article_Create(
		author,
		in.Article.Title,
		in.Article.Description,
		in.Article.Body,
		*tagList,
	), nil
}
