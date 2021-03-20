package model

import "time"

type Comment struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
	Body      string
	Author    Author
}

type Comments struct {
	Comments []Comment
}
