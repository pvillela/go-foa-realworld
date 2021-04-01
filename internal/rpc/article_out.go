package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

const dateLayout = "2006-01-02T15:04:05.999Z"

type ArticleOut struct {
	Article articleOut0
}

type articleOut0 struct {
	Slug           string   `json:"slug"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Body           *string  `json:"body"`
	TagList        []string `json:"tagList"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount int      `json:"favoritesCount"`
	Author         Profile  `json:"author"`
}

type ArticlesOut struct {
	Articles      []ArticleOut
	ArticlesCount int
}

func (s ArticleOut) FromModel(user model.User, article model.Article) ArticleOut {
	isFollowingAuthor := false
	for _, userName := range user.FollowIDs {
		if userName == article.Author.Name {
			isFollowingAuthor = true
			break
		}
	}

	favorite := false
	for _, favUser := range article.FavoritedBy {
		if user.Name == favUser.Name {
			favorite = true
			break
		}
	}

	s.Article = articleOut0{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		CreatedAt:      article.CreatedAt.UTC().Format(dateLayout),
		UpdatedAt:      article.UpdatedAt.UTC().Format(dateLayout),
		Author:         Profile{}.FromModel(article.Author, isFollowingAuthor),
		TagList:        article.TagList,
		Favorited:      favorite,
		FavoritesCount: len(article.FavoritedBy),
	}

	return s
}

func (s ArticlesOut) FromModel(user model.User, articles []model.Article) ArticlesOut {
	outs := []ArticleOut{} // return at least an empty array (not nil)

	for _, article := range articles {
		outs = append(outs, ArticleOut{}.FromModel(user, article))
	}

	s.Articles = outs
	s.ArticlesCount = len(outs)

	return s
}
