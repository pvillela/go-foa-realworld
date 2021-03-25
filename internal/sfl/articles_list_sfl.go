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
type ArticlesListSflT = func(username string, in rpc.ArticlesListIn) (*rpc.ArticlesOut, error)

func (s ArticlesListSfl) core(username string, limit, offset int, filters []model.ArticleFilter) (*model.User, model.ArticleCollection, error) {
	if limit <= 0 {
		return nil, model.ArticleCollection{}, nil
	}

	articles, err := s.ArticleGetRecentFilteredDaf(filters)
	if err != nil {
		return nil, nil, err
	}

	var user *model.User
	if username != "" {
		var err error
		user, err = s.UserGetByNameDaf(username)
		if err != nil {
			return nil, nil, err
		}
	}

	return user, model.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset), nil
}

func (s ArticlesListSfl) invoke(username string, in rpc.ArticlesListIn) (*rpc.ArticlesOut, error) {
	tagFilter := model.ArticleTagFilterOf(in.Tag)
	authorFilter := model.ArticleAuthorFilterOf(in.Author)
	favoritedFilter := model.ArticleFavoritedFilterOf(in.Favorited)
	limit := in.Limit
	offset := in.Offset

	filters := []model.ArticleFilter{tagFilter, authorFilter, favoritedFilter}

	user, articles, err := s.core(username, limit, offset, filters)
	if err != nil {
		return nil, err
	}

	articlesOut := rpc.ArticlesOutFromModel(user, articles)
	return &articlesOut, err
}

func (s ArticlesListSfl) Make() ArticlesListSflT {
	return s.invoke
}
