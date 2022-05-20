/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package model

import (
	"time"

	"github.com/pvillela/go-foa-realworld/internal/arch/util"
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

type ArticlePlus = struct {
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
	author *User,
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

	s.Slug = util.Slug(s.Title)

	return s
}
