/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleCreateIn struct {
	Article struct {
		Title       string
		Description string
		Body        *string   // mandatory
		TagList     *[]string // optional
	}
}

func (in ArticleCreateIn) ToArticle(author *model.User) model.Article {
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
	)
}
