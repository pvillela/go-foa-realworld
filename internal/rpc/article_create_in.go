/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleCreateIn struct {
	Article struct {
		Title       string
		Description string
		Body        *string
		TagList     []string // optional
	}
}

func (in ArticleCreateIn) ToArticle() model.Article {
	return model.Article{
		Slug:        fs.SlugSup(in.Article.Title),
		Title:       in.Article.Title,
		Description: in.Article.Description,
		Body:        in.Article.Body,
		TagList:     in.Article.TagList,
	}
}
