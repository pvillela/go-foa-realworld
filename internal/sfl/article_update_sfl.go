package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleUpdateSfl is the stereotype instance for the service flow that
// updates an article.
type ArticleUpdateSfl struct {
	ArticleGetAndCheckOwnerFl     fs.ArticleGetAndCheckOwnerFlT
	ArticleValidateBeforeUpdateBf fs.ArticleValidateBeforeUpdateBfT
	UserGetByNameDaf              fs.UserGetByNameDafT
	ArticleUpdateDaf              fs.ArticleUpdateDafT
	ArticleGetBySlugDaf           fs.ArticleGetBySlugDafT
	ArticleCreateDaf              fs.ArticleCreateDafT
}

// ArticleUpdateSflT is the function type instantiated by ArticleUpdateSfl.
type ArticleUpdateSflT = func(username, slug string, in rpc.ArticleUpdateIn) (*rpc.ArticleOut, error)

func (s ArticleUpdateSfl) Make() ArticleUpdateSflT {
	return func(username string, slug string, in rpc.ArticleUpdateIn) (*rpc.ArticleOut, error) {
		fieldsToUpdate := make(map[model.ArticleUpdatableField]interface{}, 3)
		if v := in.Article.Title; v != "" {
			fieldsToUpdate[model.Title] = v
		}
		if v := in.Article.Description; v != "" {
			fieldsToUpdate[model.Description] = v
		}
		if v := in.Article.Body; v != "" {
			fieldsToUpdate[model.Body] = v
		}

		article, err := s.ArticleGetAndCheckOwnerFl(slug, username)
		if err != nil {
			return nil, err
		}

		article.Update(fieldsToUpdate)
		if err := s.ArticleValidateBeforeUpdateBf(article); err != nil {
			return nil, err
		}

		user, err := s.UserGetByNameDaf(username)
		if err != nil {
			return nil, err
		}

		newSlug := fs.SlugSup(article.Title)
		article.Slug = newSlug

		var savedArticle *model.Article

		// TODO: move some of this logic to a BF
		if newSlug == slug {
			savedArticle, err = s.ArticleUpdateDaf(*article)
			if err != nil {
				return nil, err
			}
		} else {
			if _, err := s.ArticleGetBySlugDaf(newSlug); err == nil {
				return nil, fs.ErrDuplicateArticle
			}
			savedArticle, err = s.ArticleCreateDaf(*article)
			if err != nil {
				return nil, err
			}
		}

		articleOut := rpc.ArticleOut{}.FromModel(user, savedArticle)
		return &articleOut, err
	}
}
