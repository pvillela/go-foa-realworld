package fs

import "errors"

var (
	ErrDuplicateArticle = errors.New("duplicate article slug")
	ErrArticleNotFound  = errors.New("article not found")
	ErrCommentNotFound  = errors.New("comment not found")
	ErrDuplicateComment = errors.New("duplicate comment id")
	ErrUserNotFound     = errors.New("user not found")
)
