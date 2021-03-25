package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// ArticleFavoriteFl is the stereotype instance for the flow that
// designates an article as a favorite or not.
type ArticleFavoriteFl struct {
	UserGetByNameDaf    UserGetByNameDafT
	ArticleGetBySlugDaf ArticleGetBySlugDafT
	ArticleUpdateDaf    ArticleUpdateDafT
}

// ArticleFavoriteFlT is the function type instantiated by fs.ArticleFavoriteFl.
type ArticleFavoriteFlT = func(username, slug string, favorite bool) (*model.User, *model.Article, error)

func (s ArticleFavoriteFl) Make() ArticleFavoriteFlT {
	return func(username, slug string, favorite bool) (*model.User, *model.Article, error) {
		user, err := s.UserGetByNameDaf(username)
		if err != nil {
			return nil, nil, err
		}

		article, err := s.ArticleGetBySlugDaf(slug)
		if err != nil {
			return nil, nil, err
		}

		article.UpdateFavoritedBy(*user, favorite)

		updatedArticle, err := s.ArticleUpdateDaf(*article)
		if err != nil {
			return nil, nil, err
		}

		return user, updatedArticle, nil
	}
}
