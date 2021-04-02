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
type ArticleUpdateSflT = func(username, slug string, in rpc.ArticleUpdateIn) (rpc.ArticleOut, error)

func (s ArticleUpdateSfl) Make() ArticleUpdateSflT {
	return func(username string, slug string, in rpc.ArticleUpdateIn) (rpc.ArticleOut, error) {
		var zero rpc.ArticleOut

		pwArticle, err := s.ArticleGetAndCheckOwnerFl(slug, username)
		if err != nil {
			return zero, err
		}
		article := &pwArticle.Entity

		fieldsToUpdate := make(map[model.ArticleUpdatableField]interface{}, 3)
		if v := in.Article.Title; v != nil {
			fieldsToUpdate[model.Title] = *v
		}
		if v := in.Article.Description; v != nil {
			fieldsToUpdate[model.Description] = *v
		}
		if v := in.Article.Body; v != nil {
			fieldsToUpdate[model.Body] = *v
		}

		var newSlug string

		*article, newSlug = (*article).Update(fieldsToUpdate)

		if err := s.ArticleValidateBeforeUpdateBf(*article); err != nil {
			return zero, err
		}

		pwUser, err := s.UserGetByNameDaf(username)
		if err != nil {
			return zero, err
		}
		user := &pwUser.Entity

		var savedPwArticle fs.PwArticle
		savedArticle := &savedPwArticle.Entity

		// TODO: move some of this logic to a BF
		if newSlug == slug {
			savedPwArticle, err = s.ArticleUpdateDaf(pwArticle)
			if err != nil {
				return zero, err
			}
		} else {
			if _, err := s.ArticleGetBySlugDaf(newSlug); err == nil {
				return zero, fs.ErrDuplicateArticle
			}
			savedPwArticle, err = s.ArticleCreateDaf(*article)
			if err != nil {
				return zero, err
			}
		}

		articleOut := rpc.ArticleOut{}.FromModel(*user, *savedArticle)
		return articleOut, err
	}
}
