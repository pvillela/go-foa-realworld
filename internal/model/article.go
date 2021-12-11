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
	Uuid        util.Uuid
	Slug        string
	Author      User
	Title       string
	Description string
	Body        *string
	TagList     []string
	FavoritedBy []User
	Comments    []Comment
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Comment struct {
	ArticleUuid util.Uuid
	ID          int
	Author      User
	Body        *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ArticleUpdateSrc struct {
	Title       *string
	Description *string
	Body        *string
	TagList     *[]string
}

func Article_Create(
	title string,
	description string,
	body *string,
	tagList []string,
	author User,
) Article {
	now := time.Now()
	article := Article{
		Uuid:        util.NewUuid(),   // make sure this is unique index in database
		Slug:        util.Slug(title), // make sure this is unique index in database
		Author:      author,
		Title:       title,
		Description: description,
		Body:        body,
		TagList:     tagList,
		FavoritedBy: nil,
		Comments:    nil,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return article
}

func (s Article) Update(src ArticleUpdateSrc) Article {
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
	s.UpdatedAt = time.Now()

	return s
}

func (s Article) UpdateComments(comment Comment, add bool) Article {
	if add {
		s.Comments = append(s.Comments, comment)
		return s
	}

	arr := s.Comments
	extractor := func(comment Comment) int { return comment.ID }
	compValue := comment.ID
	zero := Comment{}

	// Boilerplate, repeated in next function
	index := -1
	for i := 0; i < len(arr); i++ {
		if extractor(arr[i]) == compValue {
			index = i
			break
		}
	}
	if index != -1 {
		// See https://github.com/golang/go/wiki/SliceTricks avoidance of potential memory leak.
		b := append(arr[:index], arr[index+1:]...)
		arr[len(arr)-1] = zero
		arr = b
	}

	s.Comments = arr

	return s
}

func (s Article) UpdateFavoritedBy(user User, add bool) Article {
	// This will duplicate the user if it is already in the list.
	if add {
		s.FavoritedBy = append(s.FavoritedBy, user)
		return s
	}

	for i := 0; i < len(s.FavoritedBy); i++ {
		if s.FavoritedBy[i].Name == user.Name {
			s.FavoritedBy = append(s.FavoritedBy[:i], s.FavoritedBy[i+1:]...) // memory leak ? https://github.com/golang/go/wiki/SliceTricks
		}
	}

	arr := s.FavoritedBy
	extractor := func(user User) string { return user.Name }
	compValue := user.Name
	zero := User{}

	// Boilerplate, same as in previous function
	index := -1
	for i := 0; i < len(arr); i++ {
		if extractor(arr[i]) == compValue {
			index = i
			break
		}
	}
	if index != -1 {
		// See https://github.com/golang/go/wiki/SliceTricks avoidance of potential memory leak.
		b := append(arr[:index], arr[index+1:]...)
		arr[len(arr)-1] = zero
		arr = b
	}

	s.FavoritedBy = arr

	return s

}

func Comment_Create(
	articleUuid util.Uuid,
	body *string,
	author User,
) Comment {
	now := time.Now()
	comment := Comment{
		ArticleUuid: articleUuid,
		ID:          0,
		CreatedAt:   now,
		UpdatedAt:   now,
		Body:        body,
		Author:      author,
	}
	return comment
}
