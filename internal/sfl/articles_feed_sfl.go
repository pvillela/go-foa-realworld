package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/ft"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticlesFeedSfl is the stereotype instance for the service flow that
// queries for all articles authored by other users followed by
// the current user, with optional limit and offset pagination parameters.
type ArticlesFeedSfl struct {
	UserGetByNameDaf                          ft.UserGetByNameDafT
	ArticleGetByAuthorsOrderedByMostRecentDaf ft.ArticleGetByAuthorsOrderedByMostRecentDafT
}

// FeedArticlesSflT is the function type instantiated by ArticlesFeedSfl.
type FeedArticlesSflT = func(username string, limit int, offset int) (*rpc.ArticlesOut, error)

func (s ArticlesFeedSfl) core(username string, limit, offset int) (*model.User, []model.Article, error) {
	if limit < 0 {
		return nil, []model.Article{}, nil
	}

	var user *model.User
	if username != "" {
		var errGet error
		user, errGet = s.UserGetByNameDaf(username)
		if errGet != nil {
			return nil, nil, errGet
		}
	}
	articles, err := s.ArticleGetByAuthorsOrderedByMostRecentDaf(user.FollowIDs)
	if err != nil {
		return nil, nil, err
	}

	return user, model.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset), nil
}

func (s ArticlesFeedSfl) invoke(username string, limit, offset int) (*rpc.ArticlesOut, error) {
	user, articles, err := s.core(username, limit, offset)
	if err != nil {
		return nil, err
	}
	articlesOut := rpc.ArticlesOutFromModel(user, articles)
	return &articlesOut, err
}
