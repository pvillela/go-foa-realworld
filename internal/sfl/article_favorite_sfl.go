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
type ArticleFavoriteSflT = func(username string, slug string) (rpc.ArticleOut, error)

func (s ArticleFavoriteSfl) Make() ArticleFavoriteSflT {
	return func(username string, slug string) (rpc.ArticleOut, error) {
		var zero rpc.ArticleOut
		pwUser, pwArticle, err := s.ArticleFavoriteFl(username, slug, true)
		if err != nil {
			return zero, err
		}
		articleOut := rpc.ArticleOut{}.FromModel(*pwUser.Entity(), *pwArticle.Entity())
		return articleOut, err
	}
}
