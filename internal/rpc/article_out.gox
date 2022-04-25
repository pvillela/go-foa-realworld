/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
)

const dateLayout = "2006-01-02T15:04:05.999Z"

type ArticleOut struct {
	Article articleOut0
}

type articleOut0 = struct {
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

func ArticleOut_FromModel(user model.User, article model.Article) ArticleOut {
	isFollowingAuthor := false
	for _, userName := range user.Followees {
		if userName == article.Author.Username {
			isFollowingAuthor = true
			break
		}
	}

	favorite := false
	for _, favUser := range article.FavoritedBy {
		if user.Username == favUser.Username {
			favorite = true
			break
		}
	}

	articleOut0 := articleOut0{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		CreatedAt:      article.CreatedAt.UTC().Format(dateLayout),
		UpdatedAt:      article.UpdatedAt.UTC().Format(dateLayout),
		Author:         Profile_FromModel(article.Author, isFollowingAuthor),
		TagList:        article.TagList,
		Favorited:      favorite,
		FavoritesCount: len(article.FavoritedBy),
	}

	return ArticleOut{articleOut0}
}

func ArticlesOut_FromModel(user model.User, articles []model.Article) ArticlesOut {
	outs := []ArticleOut{} // return at least an empty array (not nil)

	for _, article := range articles {
		outs = append(outs, ArticleOut_FromModel(user, article))
	}

	return ArticlesOut{outs, len(outs)}
}
