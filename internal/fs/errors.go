package fs

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	ErrDuplicateArticle     Err = "duplicate article slug"
	ErrArticleNotFound      Err = "article not found"
	ErrCommentNotFound      Err = "comment not found"
	ErrProfileNotFound      Err = "profile not found"
	ErrUserNotFound         Err = "user not found"
	ErrUnauthorizedUser     Err = "user not authorized to take this action"
	ErrAuthenticationFailed Err = "user authentication failed"
)
