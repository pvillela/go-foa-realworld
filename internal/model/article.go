package model

import "time"

type Article struct {
	Slug           string
	Title          string
	Description    string
	Body           string
	TagList        []string
	CreatedAt      time.Time
	Favorited      bool
	FavoritesCount int
	Author         Author
}

type Articles struct {
	Articles      []Article
	ArticlesCount int
}
