package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/ft"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleGetSfl is the stereotype instance for the service flow that
// retrieves an article.
type ArticleGetSfl struct {
	UserGetByNameDaf    ft.UserGetByNameDafT
	ArticleGetBySlugDaf ft.ArticleGetBySlugDafT
}

// ArticleGetSflT is the function type instantiated by ArticleGetSfl.
type ArticleGetSflT = func(username string, slug string) (*rpc.ArticleOut, error)

func (s ArticleGetSfl) core(username string, slug string) (*model.User, *model.Article, error) {
	var user *model.User
	if username != "" {
		var err error
		user, err = s.UserGetByNameDaf(username)
		if err != nil {
			return nil, nil, err
		}
	}

	article, err := s.ArticleGetBySlugDaf(slug)
	if err != nil {
		return nil, nil, err
	}

	return user, article, nil
}

func (s ArticleGetSfl) invoke(username string, slug string) (*rpc.ArticleOut, error) {
	user, article, err := s.core(username, slug)
	if err != nil {
		return nil, err
	}
	articleOut := rpc.ArticleOutFromModel(user, article)

	return &articleOut, err
}

func (s ArticleGetSfl) Make() ArticleGetSflT {
	return s.invoke
}
