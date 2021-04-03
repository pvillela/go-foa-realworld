package fs

// ArticleFavoriteFl is the stereotype instance for the flow that
// designates an article as a favorite or not.
type ArticleFavoriteFl struct {
	UserGetByNameDaf    UserGetByNameDafT
	ArticleGetBySlugDaf ArticleGetBySlugDafT
	ArticleUpdateDaf    ArticleUpdateDafT
}

// ArticleFavoriteFlT is the function type instantiated by fs.ArticleFavoriteFl.
type ArticleFavoriteFlT = func(username, slug string, favorite bool) (PwUser, PwArticle, error)

func (s ArticleFavoriteFl) Make() ArticleFavoriteFlT {
	return func(username, slug string, favorite bool) (PwUser, PwArticle, error) {
		pwUser, err := s.UserGetByNameDaf(username)
		if err != nil {
			return nil, nil, err
		}

		pwArticle, err := s.ArticleGetBySlugDaf(slug)
		if err != nil {
			return nil, nil, err
		}
		article := pwArticle.Entity()

		*article = article.UpdateFavoritedBy(*pwUser.Entity(), favorite)

		pwUpdatedArticle, err := s.ArticleUpdateDaf(pwArticle)
		if err != nil {
			return nil, nil, err
		}

		return pwUser, pwUpdatedArticle, nil
	}
}
