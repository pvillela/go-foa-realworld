/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fs

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleCheckOwnerBfT = func(article model.Article, username string) error

var ArticleCheckOwnerBfI ArticleCheckOwnerBfT = func(article model.Article, username string) error {
	if article.Author.Username != username {
		return errors.New("article not owned by user")
	}
	return nil
}
