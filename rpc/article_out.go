/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

const dateLayout = "2006-01-02T15:04:05.999Z"

type ArticleOut struct {
	Article model.ArticlePlus
}

type ArticlesOut struct {
	Articles      []ArticleOut
	ArticlesCount int
}

func ArticleOut_FromModel(articlePlus model.ArticlePlus) ArticleOut {
	return ArticleOut{articlePlus}
}

func ArticlesOut_FromModel(
	articlesPlus []model.ArticlePlus,
) ArticlesOut {
	outs := make([]ArticleOut, len(articlesPlus))

	for i := range articlesPlus {
		outs[i] = ArticleOut_FromModel(articlesPlus[i])
	}

	return ArticlesOut{outs, len(outs)}
}
