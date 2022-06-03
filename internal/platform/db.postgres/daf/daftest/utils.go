/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

func ArticlePlusesToArticles(aps []model.ArticlePlus) []model.Article {
	return util.SliceMap(aps, func(ap model.ArticlePlus) model.Article {
		return ap.ToArticle()
	})
}
