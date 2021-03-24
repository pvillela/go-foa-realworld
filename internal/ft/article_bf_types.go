package ft

import "github.com/pvillela/go-foa-realworld/internal/model"

type ArticleValidateBeforeCreateBfT = func(article model.Article) error

type ArticleValidateBeforeUpdateBfT = func(article model.Article) error

type ArticleCheckOwnerBfT = func(article model.Article, username string) error