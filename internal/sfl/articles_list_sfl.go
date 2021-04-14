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

		// TODO: move all of the contiguous logic to the DAF. The DAF should take care
		//   of the filtering logic.  All the SFL should do is pass the pointer values to the DAF.
		limit := 20
		if in.Limit != nil {
			limit = *in.Limit
			if limit < 0 {
				limit = 0
			}
		}
		offset := 0
		if in.Offset != nil {
			offset = *in.Offset
			if offset < 0 {
				offset = 0
			}
		}
		filters := make([]model.ArticleFilter, 3)
		if in.Tag != nil {
			tagFilter := model.ArticleTagFilterOf(*in.Tag)
			filters = append(filters, tagFilter)
		}
		if in.Author != nil {
			authorFilter := model.ArticleAuthorFilterOf(*in.Author)
			filters = append(filters, authorFilter)
		}
		if in.Favorited != nil {
			favoritedFilter := model.ArticleFavoritedFilterOf(*in.Favorited)
			filters = append(filters, favoritedFilter)
		}

		// TODO: Move this to the DAF
		if limit <= 0 {
			return rpc.ArticlesOut{Articles: []rpc.ArticleOut{}}, err
		}

		user := model.User{}
		if username != "" {
			user, _, err = s.UserGetByNameDaf(username)
			if err != nil {
				return zero, err
			}
		}

		// TODO: the DAF should take discrete parameters and set-up the filters internally.
		articles, err = s.ArticleGetRecentFilteredDaf(filters)
		if err != nil {
			return zero, err
		}

		articles = model.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset)

		articlesOut := rpc.ArticlesOut{}.FromModel(user, articles)

		return articlesOut, err
	}
}
