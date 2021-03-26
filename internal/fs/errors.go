package fs

import "errors"

var (
	ErrDuplicateArticle = errors.New("duplicate article slug")
	ErrArticleNotFound  = errors.New("article not found")
	ErrCommentNotFound  = errors.New("comment not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrUnauthorizedUser = errors.New("user not authorized to take this action")
)
