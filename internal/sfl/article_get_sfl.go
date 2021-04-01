package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleGetSfl is the stereotype instance for the service flow that
// retrieves an article.
type ArticleGetSfl struct {
	UserGetByNameDaf    fs.UserGetByNameDafT
	ArticleGetBySlugDaf fs.ArticleGetBySlugDafT
}

// ArticleGetSflT is the function type instantiated by ArticleGetSfl.
type ArticleGetSflT = func(username string, slug string) (rpc.ArticleOut, error)

func (s ArticleGetSfl) Make() ArticleGetSflT {
	return func(username string, slug string) (rpc.ArticleOut, error) {
		var zero rpc.ArticleOut
		var user fs.MdbUser

		if username != "" {
			var err error
			user, err = s.UserGetByNameDaf(username)
			if err != nil {
				return zero, err
			}
		}

		article, err := s.ArticleGetBySlugDaf(slug)
		if err != nil {
			return zero, err
		}

		articleOut := rpc.ArticleOut{}.FromModel(user.Entity, article.Entity)

		return articleOut, err
	}
}
