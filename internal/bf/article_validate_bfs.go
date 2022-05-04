/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package bf

import (
	"time"

	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleValidateBeforeCreateBfT = func(article model.Article) error

var ArticleValidateBeforeCreateBfI ArticleValidateBeforeCreateBfT = func(article model.Article) error {
	if article.Id == 0 || article.Slug == "" || article.Title == "" ||
		article.Description == "" || article.Body == nil {
		return ErrArticleCreateMissingFields.Make(nil)
	}
	return nil
}

type ArticleValidateBeforeUpdateBfT = func(article model.Article) error

var ArticleValidateBeforeUpdateBfI ArticleValidateBeforeUpdateBfT = func(article model.Article) error {
	if article.Id == 0 || article.Slug == "" || article.Title == "" ||
		article.Description == "" || article.Body == nil ||
		article.CreatedAt == (time.Time{}) {
		return ErrArticleUpdateMissingFields.Make(nil)
	}
	return nil
}
