package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleCreateSfl contains the dependencies required for the construction of a
// ArticleCreateSflT. It represents the creation of an article.
type ArticleCreateSfl struct {
	UserGetByNameDaf func(usename string) (*model.User, error)
	CreateArticleDaf func(article model.Article) (*model.Article, error)
}

func (s ArticleCreateSfl) core(username string, article model.Article) (*model.User, *model.Article, error) {
	user, err := s.UserGetByNameDaf(username)
	fullArticle, err := s.CreateArticleDaf(article)
	return user, fullArticle, err
}

func (s ArticleCreateSfl) Invoke(username string, in rpc.ArticleCreateIn) (*rpc.ArticleOut, error) {
	article := in.ToArticle()
	user, fullArticle, err := s.core(username, article)
	if err != nil {
		return nil, err
	}
	articleOut := rpc.ArticleOutFromModel(fullArticle, user)
	return &articleOut, err
}
