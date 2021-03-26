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
type ArticlesFeedSflT = func(username string, in rpc.ArticlesFeedIn) (*rpc.ArticlesOut, error)

func (s ArticlesFeedSfl) Make() ArticlesFeedSflT {
	return func(username string, in rpc.ArticlesFeedIn) (*rpc.ArticlesOut, error) {
		var user *model.User
		var articles []model.Article
		var err error

		limit := in.Limit
		offset := in.Offset

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

		articlesOut := rpc.ArticlesOut{}.FromModel(user, articles)

		return &articlesOut, err
	}
}
