/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package model

import (
	"time"

	"github.com/pvillela/go-foa-realworld/experimental/arch/util"
)

type Article struct {
	Id       uint
	AuthorId uint
	//Author         *User ` db:"-"`
	Title          string
	Slug           string
	Description    string
	Body           *string
	FavoritesCount int
	//FavoritedBy    []*User
	TagList []string ` db:"-"`
	//Comments  []Comment ` db:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ArticlePlus struct {
	Id             uint      `json:"-"`
	Slug           string    `json:"slug"`
	Author         Profile   `json:"author"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           *string   `json:"body"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
}

type ArticlePatch struct {
	Title       *string
	Description *string
	Body        *string
	TagList     *[]string
}

func Article_Create(
	author User,
	title string,
	description string,
	body *string,
	tagList []string,
) Article {
	article := Article{
		Slug:     util.Slug(title), // make sure this is unique index in database
		AuthorId: author.Id,
		//Author:      author,
		Title:       title,
		Description: description,
		Body:        body,
		TagList:     tagList,
	}
	return article
}

// Update returns an updated copy of the receiver.
// It does not change the slug when the title changes.
func (s Article) Update(src ArticlePatch) Article {
	if src.Title != nil {
		s.Title = *src.Title
	}
	if src.Description != nil {
		s.Description = *src.Description
	}
	if src.Body != nil {
		s.Body = src.Body
	}
	if src.TagList != nil {
		s.TagList = *src.TagList
	}

	return s
}

func (s Article) WithAdjustedFavoriteCount(delta int) Article {
	s.FavoritesCount += delta
	return s
}

func (s ArticlePlus) ToArticle() Article {
	return Article{
		Id:             s.Id,
		AuthorId:       s.Author.UserId,
		Title:          s.Title,
		Slug:           s.Slug,
		Description:    s.Description,
		Body:           s.Body,
		FavoritesCount: s.FavoritesCount,
		TagList:        s.TagList,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

func ArticlePlus_FromArticle(article Article, favorited bool, author Profile) ArticlePlus {
	return ArticlePlus{
		Id:   article.Id,
		Slug: article.Slug,
		Author: Profile{
			UserId:    author.UserId,
			Username:  author.Username,
			Bio:       author.Bio,
			Image:     author.Image,
			Following: author.Following,
		},
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        article.TagList,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		Favorited:      favorited,
		FavoritesCount: article.FavoritesCount,
	}
}
