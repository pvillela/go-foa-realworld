/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleListFl struct {
	UserGetByNameDaf                          UserGetByNameDafT
	ArticleGetByAuthorsOrderedByMostRecentDaf ArticleGetByAuthorsOrderedByMostRecentDafT
}

// ArticleListFlT is the function type instantiated by fs.ArticleListFl.
type ArticleListFlT = func(username string, limit, offset int) (model.User, []model.Article, error)

func (s ArticleListFl) Make() ArticleListFlT {
	return func(username string, limit, offset int) (model.User, []model.Article, error) {
		var zeroUser model.User

		if limit < 0 {
			return zeroUser, []model.Article{}, nil
		}

		var user model.User
		if username != "" {
			var err error
			user, _, err = s.UserGetByNameDaf(username)
			if err != nil {
				return zeroUser, nil, err
			}
		}

		articles, err := s.ArticleGetByAuthorsOrderedByMostRecentDaf(user.FollowIDs)
		if err != nil {
			return zeroUser, nil, err
		}

		return user, model.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset), nil
	}
}
