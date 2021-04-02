package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentAddSfl is the stereotype instance for the service flow that
// adds a comment to an article.
type CommentAddSfl struct {
	UserGetByNameDaf    fs.UserGetByNameDafT
	ArticleGetBySlugDaf fs.ArticleGetBySlugDafT
	CommentCreateDaf    fs.CommentCreateDafT
	ArticleUpdateDaf    fs.ArticleUpdateDafT
}

// CommentAddSflT is the function type instantiated by CommentAddSfl.
type CommentAddSflT = func(username string, in rpc.CommentAddIn) (rpc.CommentOut, error)

func (s CommentAddSfl) Make() CommentAddSflT {
	return func(username string, in rpc.CommentAddIn) (rpc.CommentOut, error) {
		var zero rpc.CommentOut
		var err error

		pwCommentAuthor, err := s.UserGetByNameDaf(username)
		if err != nil {
			return zero, err
		}
		commentAuthor := pwCommentAuthor.Entity()

		pwArticle, err := s.ArticleGetBySlugDaf(in.Slug)
		if err != nil {
			return zero, err
		}
		article := pwArticle.Entity()

		rawComment := in.ToComment(*commentAuthor)

		pwInsertedComment, err := s.CommentCreateDaf(rawComment)
		if err != nil {
			return zero, err
		}
		insertedComment := pwInsertedComment.Entity()

		article.Comments = append(article.Comments, *insertedComment)

		if _, err := s.ArticleUpdateDaf(pwArticle); err != nil {
			return zero, err
		}

		commentOut := rpc.CommentOut{}.FromModel(*insertedComment)
		return commentOut, err
	}
}
