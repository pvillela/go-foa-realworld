package fs

// ArticleGetAndCheckOwnerFl is the stereotype instance for the flow that
// checks if a given article's author's username matches a given username.
type ArticleGetAndCheckOwnerFl struct {
	ArticleGetBySlugDaf ArticleGetBySlugDafT
	ArticleCheckOwnerBf ArticleCheckOwnerBfT
}

// ArticleGetAndCheckOwnerFlT is the function type instantiated by fs.ArticleGetAndCheckOwnerFl.
type ArticleGetAndCheckOwnerFlT = func(username, slug string) (PwArticle, error)

func (s ArticleGetAndCheckOwnerFl) Make() ArticleGetAndCheckOwnerFlT {
	return func(slug string, username string) (PwArticle, error) {
		pwArticle, err := s.ArticleGetBySlugDaf(slug)
		if err != nil {
			return nil, err
		}

		if err := s.ArticleCheckOwnerBf(*pwArticle.Entity(), username); err != nil {
			return nil, err
		}

		return pwArticle, err
	}
}
