/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package model

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"time"
)

type Article struct {
	Slug        string
	Title       string
	Description string
	Body        *string
	TagList     []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FavoritedBy []User
	Author      User
	Comments    []Comment
}

type Comment struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      *string
	Author    User
}

type ArticleUpdatableField int

const (
	Title ArticleUpdatableField = iota
	Description
	Body
	TagList
)

func (s Article) Update(fieldsToUpdate map[ArticleUpdatableField]interface{}) (Article, string) {
	article := s
	for k, v := range fieldsToUpdate {
		switch k {
		case Title:
			article.Title = v.(string)
		case Description:
			article.Description = v.(string)
		case Body:
			article.Body = v.(*string)
		case TagList:
			article.TagList = v.([]string)
		}
	}
	newSlug := fs.SlugSup(article.Title)
	article.Slug = newSlug
	return article, newSlug
}

//TODO: move to BF -- article_filter_bfs.go
type ArticleFilter = func(Article) bool

//TODO: move to BF -- article_filter_bfs.go
func ArticleTagFilterOf(tag string) ArticleFilter {
	return func(article Article) bool {
		for _, articleTag := range article.TagList {
			if articleTag == tag {
				return true
			}
		}
		return false
	}
}

//TODO: move to BF -- article_filter_bfs.go
func ArticleAuthorFilterOf(authorName string) ArticleFilter {
	return func(article Article) bool {
		return article.Author.Name == authorName
	}
}

//TODO: move to BF -- article_filter_bfs.go
func ArticleFavoritedFilterOf(username string) ArticleFilter {
	return func(article Article) bool {
		if username == "" {
			return false
		}
		for _, user := range article.FavoritedBy {
			if user.Name == username {
				return true
			}
		}
		return false
	}
}

type ArticleCollection []Article

func (articles ArticleCollection) ApplyLimitAndOffset(limit, offset int) ArticleCollection {
	if limit <= 0 {
		return []Article{}
	}

	articlesSize := len(articles)
	min := offset
	if min < 0 {
		min = 0
	}

	if min > articlesSize {
		return []Article{}
	}

	max := min + limit
	if max > articlesSize {
		max = articlesSize
	}

	return articles[min:max]
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
