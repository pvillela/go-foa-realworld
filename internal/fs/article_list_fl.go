package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleListFl struct {
	UserGetByNameDaf                          UserGetByNameDafT
	ArticleGetByAuthorsOrderedByMostRecentDaf ArticleGetByAuthorsOrderedByMostRecentDafT
}

// ArticleListFlT is the function type instantiated by fs.ArticleListFl.
type ArticleListFlT = func(username string, limit, offset int) (*model.User, []model.Article, error)

func (s ArticleListFl) invoke(username string, limit, offset int) (*model.User, []model.Article, error) {
	if limit < 0 {
		return nil, []model.Article{}, nil
	}

	var user *model.User
	if username != "" {
		var err error
		user, err = s.UserGetByNameDaf(username)
		if err != nil {
			return nil, nil, err
		}
	}
	articles, err := s.ArticleGetByAuthorsOrderedByMostRecentDaf(user.FollowIDs)
	if err != nil {
		return nil, nil, err
	}

	return user, model.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset), nil
}

func (s ArticleListFl) Make() ArticleListFlT {
	return s.invoke
}
