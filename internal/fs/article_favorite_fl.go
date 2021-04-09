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
		user, rcUser, err := s.UserGetByNameDaf(username)
		if err != nil {
			return PwUser{}, PwArticle{}, err
		}

		article, rcArticle, err := s.ArticleGetBySlugDaf(slug)
		if err != nil {
			return PwUser{}, PwArticle{}, err
		}

		article = article.UpdateFavoritedBy(user, favorite)

		article, rcArticle, err = s.ArticleUpdateDaf(article, rcArticle)
		if err != nil {
			return PwUser{}, PwArticle{}, err
		}

		return PwUser{rcUser, user}, PwArticle{rcArticle, article}, nil
	}
}
