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
	Author         Profile  `json:"author"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Body           *string  `json:"body"`
	TagList        []string `json:"tagList"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount int      `json:"favoritesCount"`
}

type ArticlesOut struct {
	Articles      []ArticleOut
	ArticlesCount int
}

func ArticleOut_FromModel(article model.Article, followsAuthor bool, likesArticle bool) ArticleOut {
	articleOut0 := articleOut0{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		CreatedAt:      article.CreatedAt.UTC().Format(dateLayout),
		UpdatedAt:      article.UpdatedAt.UTC().Format(dateLayout),
		Author:         Profile_FromModel(article.Author, followsAuthor),
		TagList:        article.TagList,
		Favorited:      likesArticle,
		FavoritesCount: article.FavoritesCount,
	}
	return ArticleOut{articleOut0}
}

func ArticlesOut_FromModel(
	articles []model.Article,
	followsAuthors []bool,
	likesArticles []bool,
) ArticlesOut {
	outs := []ArticleOut{} // return at least an empty array (not nil)

	for i, article := range articles {
		outs = append(outs, ArticleOut_FromModel(article, followsAuthors[i], likesArticles[i]))
	}

	return ArticlesOut{outs, len(outs)}
}
