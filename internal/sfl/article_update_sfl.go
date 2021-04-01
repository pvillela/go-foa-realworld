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

		mdbArticle, err := s.ArticleGetAndCheckOwnerFl(slug, username)
		if err != nil {
			return zero, err
		}
		article := &mdbArticle.Entity

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

		mdbUser, err := s.UserGetByNameDaf(username)
		if err != nil {
			return zero, err
		}
		user := &mdbUser.Entity

		var savedMdbArticle fs.MdbArticle
		savedArticle := &savedMdbArticle.Entity

		// TODO: move some of this logic to a BF
		if newSlug == slug {
			savedMdbArticle, err = s.ArticleUpdateDaf(mdbArticle)
			if err != nil {
				return zero, err
			}
		} else {
			if _, err := s.ArticleGetBySlugDaf(newSlug); err == nil {
				return zero, fs.ErrDuplicateArticle
			}
			savedMdbArticle, err = s.ArticleCreateDaf(*article)
			if err != nil {
				return zero, err
			}
		}

		articleOut := rpc.ArticleOut{}.FromModel(*user, *savedArticle)
		return articleOut, err
	}
}
