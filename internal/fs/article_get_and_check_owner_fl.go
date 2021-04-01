package fs

// ArticleGetAndCheckOwnerFl is the stereotype instance for the flow that
// checks if a given article's author's username matches a given username.
type ArticleGetAndCheckOwnerFl struct {
	ArticleGetBySlugDaf ArticleGetBySlugDafT
	ArticleCheckOwnerBf ArticleCheckOwnerBfT
}

// ArticleGetAndCheckOwnerFlT is the function type instantiated by fs.ArticleGetAndCheckOwnerFl.
type ArticleGetAndCheckOwnerFlT = func(username, slug string) (MdbArticle, error)

func (s ArticleGetAndCheckOwnerFl) Make() ArticleGetAndCheckOwnerFlT {
	return func(slug string, username string) (MdbArticle, error) {
		var zero MdbArticle

		article, err := s.ArticleGetBySlugDaf(slug)
		if err != nil {
			return zero, err
		}

		if err := s.ArticleCheckOwnerBf(article.Entity, username); err != nil {
			return zero, err
		}

		return article, err
	}
}
