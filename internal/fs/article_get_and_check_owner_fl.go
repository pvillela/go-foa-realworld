/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// ArticleGetAndCheckOwnerFl is the stereotype instance for the flow that
// checks if a given article's author's username matches a given username.
type ArticleGetAndCheckOwnerFl struct {
	ArticleGetBySlugDaf ArticleGetBySlugDafT
	ArticleCheckOwnerBf ArticleCheckOwnerBfT
}

// ArticleGetAndCheckOwnerFlT is the function type instantiated by fs.ArticleGetAndCheckOwnerFl.
type ArticleGetAndCheckOwnerFlT = func(username, slug string) (model.Article, db.RecCtx, error)

func (s ArticleGetAndCheckOwnerFl) Make() ArticleGetAndCheckOwnerFlT {
	return func(slug string, username string) (model.Article, db.RecCtx, error) {
		article, rc, err := s.ArticleGetBySlugDaf(slug)
		if err != nil {
			return model.Article{}, nil, err
		}

		if err := s.ArticleCheckOwnerBf(article, username); err != nil {
			return model.Article{}, nil, err
		}

		return article, rc, err
	}
}
