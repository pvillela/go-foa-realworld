/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleUnfavoriteSfl is the stereotype instance for the service flow that
// designates an article as a favorite.
type ArticleUnfavoriteSfl struct {
	ArticleFavoriteFl fs.ArticleFavoriteFlT
}

// ArticleUnfavoriteSflT is the function type instantiated by ArticleUnfavoriteSfl.
type ArticleUnfavoriteSflT = func(username string, slug string) (rpc.ArticleOut, error)

func (s ArticleUnfavoriteSfl) Make() ArticleUnfavoriteSflT {
	return func(username string, slug string) (rpc.ArticleOut, error) {
		var zero rpc.ArticleOut
		pwUser, pwArticle, err := s.ArticleFavoriteFl(username, slug, false)
		if err != nil {
			return zero, err
		}
		articleOut := rpc.ArticleOut{}.FromModel(pwUser.Entity, pwArticle.Entity)
		return articleOut, err
	}
}
