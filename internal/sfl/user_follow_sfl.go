package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentAddSfl is the stereotype instance for the service flow that
// causes the current user start following a given other user.
type UserFollowSfl struct {
	UserGetByNameDaf fs.UserGetByNameDafT
	UserUpdateDaf fs.UserUpdateDafT
}

// UserFollowSflT is the function type instantiated by UserFollowSfl.
type UserFollowSflT = func(username string, followedUsername string) (*rpc.ProfileOut, error)

func (s UserFollowSfl) Make() UserFollowSflT {
	return func(username string, followedUsername string) (*rpc.ProfileOut, error) {
		panic("todo")
	}
}

func (s UserFollowSfl) core(userName, followedUsername string, follow bool) (*model.User, error) {
	user, err := i.userRW.GetByName(userName)
	if err != nil {
		return nil, err
	}
	if user.Name != userName {
		return nil, errWrongUser
	}
	if user == nil {
		return nil, ErrNotFound
	}

	user.UpdateFollowees(followedUsername, follow)

	if err := i.userRW.Save(*user); err != nil {
		return nil, err
	}

	return user, nil
}
