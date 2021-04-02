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
			pwUser, err := s.UserGetByNameDaf(username)
			if err != nil {
				return zeroUser, nil, err
			}
			user = pwUser.Entity
		}

		pwArticles, err := s.ArticleGetByAuthorsOrderedByMostRecentDaf(user.FollowIDs)
		if err != nil {
			return zeroUser, nil, err
		}

		articles := make([]model.Article, len(pwArticles))
		for i := 0; i < len(pwArticles); i++ {
			articles[i] = pwArticles[i].Entity
		}

		return user, model.ArticleCollection(articles).ApplyLimitAndOffset(limit, offset), nil
	}
}
