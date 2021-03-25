package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticlesFeedSfl is the stereotype instance for the service flow that
// queries for all articles authored by other users followed by
// the current user, with optional limit and offset pagination parameters.
type ArticlesFeedSfl struct {
	UserGetByNameDaf                          fs.UserGetByNameDafT
	ArticleGetByAuthorsOrderedByMostRecentDaf fs.ArticleGetByAuthorsOrderedByMostRecentDafT
}

// ArticlesFeedSflT is the function type instantiated by ArticlesFeedSfl.
type ArticlesFeedSflT = func(username string, limit int, offset int) (*rpc.ArticlesOut, error)

func (s ArticlesFeedSfl) Make() ArticlesFeedSflT {
	return func(username string, limit, offset int) (*rpc.ArticlesOut, error) {
		var user *model.User
		var articles []model.Article
		var err error

		if limit > 0 {
			if username != "" {
				user, err = s.UserGetByNameDaf(username)
				if err != nil {
					return nil, err
				}
			}

			articles, err = s.ArticleGetByAuthorsOrderedByMostRecentDaf(user.FollowIDs)
			if err != nil {
				return nil, err
			}
		}

		articles = model.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset)

		articlesOut := rpc.ArticlesOutFromModel(user, articles)

		return &articlesOut, err
	}
}
