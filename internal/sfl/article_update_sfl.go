package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleUpdateSfl is the stereotype instance for the service flow that
// updates an article.
type ArticleUpdateSfl struct {
	GetArticleAndCheckOwnerFl     func(slug string, username string) (*model.Article, error)
	ArticleValidateBeforeUpdateBf func(model.Article) error
	UserGetByNameDaf              func(userName string) (*model.User, error)
	ArticleUpdateDaf              func(article model.Article) (*model.Article, error)
	ArticleGetBySlugDaf           func(slug string) (*model.Article, error)
	ArticleCreateDaf              func(article model.Article) (*model.Article, error)
}

// ArticleUpdateSflT is the function type instantiated by ArticleUpdateSfl.
type ArticleUpdateSflT = func(username, slug string, in rpc.ArticleUpdateIn) (*rpc.ArticleOut, error)

func (s ArticleUpdateSfl) core(
	username string,
	slug string,
	fieldsToUpdate map[model.ArticleUpdatableField]interface{},
) (*model.User, *model.Article, error) {
	article, err := s.GetArticleAndCheckOwnerFl(slug, username)
	if err != nil {
		return nil, nil, err
	}

	article.Update(fieldsToUpdate)

	if err := s.ArticleValidateBeforeUpdateBf(*article); err != nil {
		return nil, nil, err
	}

	user, err := s.UserGetByNameDaf(username)
	if err != nil {
		return nil, nil, err
	}

	newSlug := fs.SlugSup(article.Title)
	article.Slug = newSlug

	var savedArticle *model.Article

	if newSlug == slug {
		savedArticle, err = s.ArticleUpdateDaf(*article)
		if err != nil {
			return nil, nil, err
		}
	} else {
		if _, err := s.ArticleGetBySlugDaf(newSlug); err == nil {
			return nil, nil, fs.ErrDuplicateArticle
		}
		savedArticle, err = s.ArticleCreateDaf(*article)
		if err != nil {
			return nil, nil, err
		}
	}

	return user, savedArticle, nil
}

func (s ArticleUpdateSfl) invoke(username string, slug string, in rpc.ArticleUpdateIn) (*rpc.ArticleOut, error) {
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

	user, article, err := s.core(username, slug, fieldsToUpdate)
	if err != nil {
		return nil, err
	}
	articleOut := rpc.ArticleOutFromModel(user, article)
	return &articleOut, err
}

func (s ArticleUpdateSfl) Make() ArticleUpdateSflT {
	return s.invoke
}
