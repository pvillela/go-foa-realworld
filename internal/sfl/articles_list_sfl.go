/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticlesListSfl is the stereotype instance for the service flow that
// retrieve recent articles based on a set of query parameters.
type ArticlesListSfl struct {
	UserGetByNameDaf            fs.UserGetByNameDafT
	ArticleGetRecentFilteredDaf fs.ArticleGetRecentFilteredDafT
}

// ArticlesListSflT is the function type instantiated by ArticlesListSfl.
type ArticlesListSflT = func(username string, in rpc.ArticlesListIn) (rpc.ArticlesOut, error)

func (s ArticlesListSfl) Make() ArticlesListSflT {
	return func(username string, in rpc.ArticlesListIn) (rpc.ArticlesOut, error) {
		var zero rpc.ArticlesOut
		var articles []model.Article
		var err error

		limit := in.Limit
		offset := in.Offset

		tagFilter := model.ArticleTagFilterOf(in.Tag)
		authorFilter := model.ArticleAuthorFilterOf(in.Author)
		favoritedFilter := model.ArticleFavoritedFilterOf(in.Favorited)
		filters := []model.ArticleFilter{tagFilter, authorFilter, favoritedFilter}

		if limit <= 0 {
			return rpc.ArticlesOut{
				Articles: []rpc.ArticleOut{},
			}, err
		}

		user := model.User{}
		if username != "" {
			user, _, err = s.UserGetByNameDaf(username)
			if err != nil {
				return zero, err
			}
		}

		articles, err = s.ArticleGetRecentFilteredDaf(filters)
		if err != nil {
			return zero, err
		}

		articles = model.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset)

		articlesOut := rpc.ArticlesOut{}.FromModel(user, articles)

		return articlesOut, err
	}
}
