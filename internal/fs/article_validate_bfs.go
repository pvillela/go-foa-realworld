/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleValidateBeforeCreateBf struct{}

type ArticleValidateBeforeCreateBfT = func(article model.Article) error

func (s ArticleValidateBeforeCreateBf) Make() ArticleValidateBeforeCreateBfT {
	return func(article model.Article) error {
		return nil
	}
}

type ArticleValidateBeforeUpdateBf struct{}

type ArticleValidateBeforeUpdateBfT = func(article model.Article) error

func (s ArticleValidateBeforeUpdateBf) Make() ArticleValidateBeforeUpdateBfT {
	return func(article model.Article) error {
		return nil
	}
}
