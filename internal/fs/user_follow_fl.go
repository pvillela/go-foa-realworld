/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

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
type UserFollowFlT = func(username string, followedUsername string, follow bool, txn db.Txn) (model.User, db.RecCtx, error)

func (s UserFollowFl) Make() UserFollowFlT {
	return func(username string, followedUsername string, follow bool, txn db.Txn) (model.User, db.RecCtx, error) {
		user, rc, err := s.UserGetByNameDaf(username)
		if err != nil {
			return model.User{}, nil, err
		}

		user = user.UpdateFollowees(followedUsername, follow)

		if rc, err = s.UserUpdateDaf(user, rc, txn); err != nil {
			return model.User{}, nil, err
		}

		return user, rc, nil
	}
}
