package fs

// ArticleFavoriteFl is the stereotype instance for the flow that
// designates an article as a favorite or not.
type ArticleFavoriteFl struct {
	UserGetByNameDaf    UserGetByNameDafT
	ArticleGetBySlugDaf ArticleGetBySlugDafT
	ArticleUpdateDaf    ArticleUpdateDafT
}

// ArticleFavoriteFlT is the function type instantiated by fs.ArticleFavoriteFl.
type ArticleFavoriteFlT = func(username, slug string, favorite bool) (MdbUser, MdbArticle, error)

func (s ArticleFavoriteFl) Make() ArticleFavoriteFlT {
	return func(username, slug string, favorite bool) (MdbUser, MdbArticle, error) {
		var zeroUser MdbUser
		var zeroArticle MdbArticle

		user, err := s.UserGetByNameDaf(username)
		if err != nil {
			return zeroUser, zeroArticle, err
		}

		article, err := s.ArticleGetBySlugDaf(slug)
		if err != nil {
			return zeroUser, zeroArticle, err
		}

		article.UpdateFavoritedBy(user.Entity, favorite)

		updatedArticle, err := s.ArticleUpdateDaf(article)
		if err != nil {
			return zeroUser, zeroArticle, err
		}

		return user, updatedArticle, nil
	}
}
