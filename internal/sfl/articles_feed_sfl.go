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
type ArticlesFeedSflT = func(username string, in rpc.ArticlesFeedIn) (rpc.ArticlesOut, error)

func (s ArticlesFeedSfl) Make() ArticlesFeedSflT {
	return func(username string, in rpc.ArticlesFeedIn) (rpc.ArticlesOut, error) {
		var zero rpc.ArticlesOut
		var pwUser fs.PwUser
		var user *model.User
		var articles []model.Article
		var err error

		if username == "" {
			return zero, fs.ErrAuthenticationFailed
		}

		limit := in.Limit
		offset := in.Offset

		if limit <= 0 {
			return rpc.ArticlesOut{
				Articles: []rpc.ArticleOut{},
			}, err
		}

		if username != "" {
			pwUser, err = s.UserGetByNameDaf(username)
			if err != nil {
				return zero, err
			}
		}

		pwUser, err = s.UserGetByNameDaf(username)
		if err != nil {
			return zero, err
		}
		user = pwUser.Entity()

		articles, err = s.ArticleGetByAuthorsOrderedByMostRecentDaf(user.FollowIDs)
		if err != nil {
			return zero, err
		}

		articles = model.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset)

		articlesOut := rpc.ArticlesOut{}.FromModel(*user, articles)

		return articlesOut, err
	}
}
