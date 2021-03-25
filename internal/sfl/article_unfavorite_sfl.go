package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleUnfavoriteSfl is the stereotype instance for the service flow that
// designates an article as a favorite.
type ArticleUnfavoriteSfl struct {
	ArticleFavoriteFl fs.ArticleFavoriteFlT
}

// ArticleUnfavoriteSflT is the function type instantiated by ArticleUnfavoriteSfl.
type ArticleUnfavoriteSflT = func(username string, slug string) (*rpc.ArticleOut, error)

func (s ArticleUnfavoriteSfl) Make() ArticleUnfavoriteSflT {
	return func(username string, slug string) (*rpc.ArticleOut, error) {
		user, article, err := s.ArticleFavoriteFl(username, slug, false)
		if err != nil {
			return nil, err
		}
		articleOut := rpc.ArticleOutFromModel(user, article)
		return &articleOut, err
	}
}
