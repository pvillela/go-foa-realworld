package fn

import "errors"

var (
	ErrDuplicateArticle = errors.New("duplicate article slug")
	ErrArticleNotFound  = errors.New("article not found")
	ErrUserNotFound     = errors.New("user not found")
)
