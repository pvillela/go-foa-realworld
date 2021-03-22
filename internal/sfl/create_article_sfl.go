package sfl

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/fn"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CreateArticleSflPre contains the dependencies required for the construction of a
// CreateArticleSflS. It represents the creation of an article.
type CreateArticleSflPre struct {
	CreateArticleDaf func(article model.Article) (*model.Article, error)
}

type CreateArticleSflS struct {
	CreateArticleSflPre
	Other struct{}
}

func (pre CreateArticleSflPre) Prep() CreateArticleSflS {
	res := CreateArticleSflS{
		CreateArticleSflPre: pre,
	}
	return res
}

func (s CreateArticleSflS) Core(username string, article model.Article) (*model.Article, error) {
	return s.CreateArticleDaf(article)
}

type CreateArticleSfl = func(username string, articleIn rpc.CreateArticleIn) model.Article

func (s CreateArticleSflS) Sfl(username string, articleIn rpc.CreateArticleIn) (*model.Article, error) {
	article := ArticleFromReq(articleIn)
	pArticle, err := s.Core(username, article)
	return pArticle, err
}

func ArticleFromReq(in rpc.CreateArticleIn) model.Article {
	return model.Article{}
}

func foo() {
	x := CreateArticleSflPre{
		CreateArticleDaf: fn.ArticleDafS{}.Create,
	}
	fmt.Println(x)
}
