package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/fn"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

type ArticleCreateIn struct {
	Article struct {
		Title       string
		Description string
		Body        string
		TagList     []string // optional
	}
}

func (in *ArticleCreateIn) ToArticle() model.Article {
	return model.Article{
		Slug:        fn.SlugSup((*in).Article.Title),
		Title:       (*in).Article.Title,
		Description: (*in).Article.Description,
		Body:        (*in).Article.Body,
		TagList:     (*in).Article.TagList,
	}
}
