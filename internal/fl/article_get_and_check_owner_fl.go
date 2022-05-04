/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fl

import (
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

// ArticleGetAndCheckOwnerFlT is the type of the stereotype instance for the flow that
// checks if a given article's author's username matches a given username, with hard-wired BF dependencies.
type ArticleGetAndCheckOwnerFlT = func(username, slug string) (model.Article, RecCtxArticle, error)

// ArticleGetAndCheckOwnerFlC is the function that constructs a stereotype instance of type
// ArticleGetAndCheckOwnerFlT with hard-wired BF dependencies.
func ArticleGetAndCheckOwnerFlC(
	articleGetBySlugDaf daf.ArticleGetBySlugDafT,
) ArticleGetAndCheckOwnerFlT {
	articleCheckOwnerBf := bf.ArticleCheckOwnerBfI
	return ArticleGetAndCheckOwnerFlC0(
		articleGetBySlugDaf,
		articleCheckOwnerBf,
	)
}

// ArticleGetAndCheckOwnerFlC0 is the function that constructs a stereotype instance of type
// ArticleGetAndCheckOwnerFlT without hard-wired BF dependencies.
func ArticleGetAndCheckOwnerFlC0(
	articleGetBySlugDaf daf.ArticleGetBySlugDafT,
	articleCheckOwnerBf bf.ArticleCheckOwnerBfT,
) ArticleGetAndCheckOwnerFlT {
	return func(slug string, username string) (model.Article, RecCtxArticle, error) {
		article, rc, err := articleGetBySlugDaf(slug)
		if err != nil {
			return model.Article{}, RecCtxArticle{}, err
		}

		if err := articleCheckOwnerBf(article, username); err != nil {
			return model.Article{}, RecCtxArticle{}, err
		}

		return article, rc, err
	}
}
