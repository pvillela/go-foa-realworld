package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fn"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleUpdateSflS contains the dependencies required for the construction of a
// ArticleUpdateSfl. It represents the updating of an article.
type ArticleUpdateSflS struct {
	ArticleGetBySlugDaf           func(slug string) (*model.Article, error)
	ArticleValidateBeforeUpdateBf func(model.Article) error
	UserGetByNameDaf              func(userName string) (*model.User, error)
	SlugBf                        func(string) string
	ArticleUpdateDaf              func(article model.Article) (*model.Article, error)
	ArticleCreateDaf              func(article model.Article) (*model.Article, error)
}

// ArticleUpdateSfl is the type of a function that takes an rpc.ArticleUpdateIn as input and
// returns a model.Article.
type ArticleUpdateSfl = func(articleIn rpc.ArticleUpdateIn) model.Article

func (s ArticleUpdateSflS) core(
	username string,
	slug string,
	fieldsToUpdate map[model.ArticleUpdatableField]interface{},
) (*model.User, *model.Article, error) {
	article, err := s.ArticleGetBySlugDaf(slug)
	if err != nil {
		return nil, nil, err
	}

	if err := fn.CheckArticleUserOwnershipBf(*article, username); err != nil {
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

	newSlug := s.SlugBf(article.Title)
	article.Slug = newSlug

	var savedArticle *model.Article

	if newSlug == slug {
		savedArticle, err = s.ArticleUpdateDaf(*article)
		if err != nil {
			return nil, nil, err
		}
	} else {
		if _, err := s.ArticleGetBySlugDaf(newSlug); err == nil {
			return nil, nil, fn.ErrDuplicateArticle
		}
		savedArticle, err = s.ArticleCreateDaf(*article)
		if err != nil {
			return nil, nil, err
		}
	}

	return user, savedArticle, nil
}

func (s ArticleUpdateSflS) Invoke(username string, slug string, in rpc.ArticleUpdateIn) (rpc.ArticleOut, error) {
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
	articleOut := rpc.ArticleOut{}
	if err != nil {
		return articleOut, err
	}
	articleOut = rpc.ArticleOutFromModel(article, user)
	return articleOut, err
}
