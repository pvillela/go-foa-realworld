/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package bf

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleCheckOwnerBfT = func(article model.ArticlePlus, username string) error

var ArticleCheckOwnerBfI ArticleCheckOwnerBfT = func(article model.ArticlePlus, username string) error {
	if article.Author.Username != username {
		return ErrUnauthorizedUser.Make(nil, username)
	}
	return nil
}
