package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleFavoriteSfl is the stereotype instance for the service flow that
// designates an article as a favorite.
type ArticleFavoriteSfl struct {
	ArticleFavoriteFl fs.ArticleFavoriteFlT
}

// ArticleFavoriteSflT is the function type instantiated by ArticleFavoriteSfl.
type ArticleFavoriteSflT = func(username string, slug string) (*rpc.ArticleOut, error)

func (s ArticleFavoriteSfl) invoke(username string, slug string) (*rpc.ArticleOut, error) {
	user, article, err := s.ArticleFavoriteFl(username, slug, true)
	if err != nil {
		return nil, err
	}
	articleOut := rpc.ArticleOutFromModel(user, article)
	return &articleOut, err
}

func (s ArticleFavoriteSfl) Make() ArticleFavoriteSflT {
	return s.invoke
}
