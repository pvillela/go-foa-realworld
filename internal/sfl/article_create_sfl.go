package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleCreateSfl is the stereotype instance for the service flow that
// creates an article.
type ArticleCreateSfl struct {
	UserGetByNameDaf func(usename string) (*model.User, error)
	CreateArticleDaf func(article model.Article) (*model.Article, error)
}

func (s ArticleCreateSfl) core(username string, article model.Article) (*model.User, *model.Article, error) {
	user, err := s.UserGetByNameDaf(username)
	fullArticle, err := s.CreateArticleDaf(article)
	return user, fullArticle, err
}

// ArticleCreateSflT is the function type instantiated by ArticleCreateSfl.
type ArticleCreateSflT = func(username string, in rpc.ArticleCreateIn) (*rpc.ArticleOut, error)

func (s ArticleCreateSfl) invoke(username string, in rpc.ArticleCreateIn) (*rpc.ArticleOut, error) {
	article := in.ToArticle()
	user, fullArticle, err := s.core(username, article)
	if err != nil {
		return nil, err
	}
	articleOut := rpc.ArticleOutFromModel(user, fullArticle)
	return &articleOut, err
}

func (s ArticleCreateSfl) Make() ArticleCreateSflT {
	return s.invoke
}
