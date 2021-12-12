/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// ArticleGetAndCheckOwnerFlT is the type of the stereotype instance for the flow that
// checks if a given article's author's username matches a given username.
type ArticleGetAndCheckOwnerFlT = func(username, slug string) (model.Article, RecCtxArticle, error)

// ArticleGetAndCheckOwnerFlC is the function that constructs a stereotype instance of type
// ArticleGetAndCheckOwnerFlT.
func ArticleGetAndCheckOwnerFlC(
	articleGetBySlugDaf ArticleGetBySlugDafT,
) ArticleGetAndCheckOwnerFlT {
	articleCheckOwnerBf := ArticleCheckOwnerBfI
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
