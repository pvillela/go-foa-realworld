package rpc

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
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
		Slug:        fs.SlugSup((*in).Article.Title),
		Title:       (*in).Article.Title,
		Description: (*in).Article.Description,
		Body:        (*in).Article.Body,
		TagList:     (*in).Article.TagList,
	}
}
