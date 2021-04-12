/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
)

// ArticleFavoriteSfl is the stereotype instance for the service flow that
// deletes an article.
type ArticleDeleteSfl struct {
	ArticleGetAndCheckOwnerFl fs.ArticleGetAndCheckOwnerFlT
	ArticleDeleteDaf          fs.ArticleDeleteDafT
}

// ArticleDeleteSflT is the function type instantiated by ArticleDeleteSfl.
type ArticleDeleteSflT = func(username, slug string) error

func (s ArticleDeleteSfl) Make() ArticleDeleteSflT {
	return func(username string, slug string) error {
		_, _, err := s.ArticleGetAndCheckOwnerFl(username, slug)
		if err != nil {
			return err
		}

		return s.ArticleDeleteDaf(slug)
	}
}
