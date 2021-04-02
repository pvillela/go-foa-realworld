package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleCreateSfl is the stereotype instance for the service flow that
// creates an article.
type ArticleCreateSfl struct {
	UserGetByNameDaf              fs.UserGetByNameDafT
	ArticleValidateBeforeCreateBf fs.ArticleValidateBeforeCreateBfT
	ArticleCreateDaf              fs.ArticleCreateDafT
	TagAddDaf                     fs.TagAddDafT
}

// ArticleCreateSflT is the function type instantiated by ArticleCreateSfl.
type ArticleCreateSflT = func(username string, in rpc.ArticleCreateIn) (rpc.ArticleOut, error)

func (s ArticleCreateSfl) Make() ArticleCreateSflT {
	return func(username string, in rpc.ArticleCreateIn) (rpc.ArticleOut, error) {
		zero := rpc.ArticleOut{}
		article := in.ToArticle()

		pwUser, err := s.UserGetByNameDaf(username)
		if err != nil {
			return zero, err
		}

		if err := s.ArticleValidateBeforeCreateBf(article); err != nil {
			return zero, err
		}

		pwArticle, err := s.ArticleCreateDaf(article)
		if err != nil {
			return zero, err
		}

		if err := s.TagAddDaf(article.TagList); err != nil {
			return zero, err
		}

		articleOut := rpc.ArticleOut{}.FromModel(pwUser.Entity, pwArticle.Entity)
		return articleOut, err
	}
}
