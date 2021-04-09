package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// CommentAddSfl is the stereotype instance for the service flow that
// causes the current user start following a given other user.
type UserFollowFl struct {
	UserGetByNameDaf UserGetByNameDafT
	UserUpdateDaf    UserUpdateDafT
}

// UserFollowFlT is the function type instantiated by UserFollowFl.
type UserFollowFlT = func(username string, followedUsername string, follow bool) (model.User, db.RecCtx, error)

func (s UserFollowFl) Make() UserFollowFlT {
	return func(username string, followedUsername string, follow bool) (model.User, db.RecCtx, error) {
		user, rc, err := s.UserGetByNameDaf(username)
		if err != nil {
			return model.User{}, nil, err
		}

		user = user.UpdateFollowees(followedUsername, follow)

		if user, rc, err = s.UserUpdateDaf(user, rc); err != nil {
			return model.User{}, nil, err
		}

		return user, rc, nil
	}
}
