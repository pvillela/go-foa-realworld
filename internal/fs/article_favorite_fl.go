package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/ft"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// ArticleFavoriteFl is the stereotype instance for the flow that
// designates an article as a favorite or not.
type ArticleFavoriteFl struct {
	UserGetByNameDaf    ft.UserGetByNameDafT
	ArticleGetBySlugDaf ft.ArticleGetBySlugDafT
	ArticleUpdateDaf    ft.ArticleUpdateDafT
}

func (s ArticleFavoriteFl) invoke(username, slug string, favorite bool) (*model.User, *model.Article, error) {
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

func (s ArticleFavoriteFl) Make() ft.ArticleFavoriteFlT {
	return s.invoke
}
