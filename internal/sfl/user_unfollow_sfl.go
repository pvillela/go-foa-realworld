/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentAddSfl is the stereotype instance for the service flow that
// causes the current user start following a given other user.
type UserUnfollowSfl struct {
	UserFollowFl fs.UserFollowFlT
}

// UserUnfollowSflT is the function type instantiated by UserUnfollowSfl.
type UserUnfollowSflT = func(username string, followedUsername string) (rpc.ProfileOut, error)

func (s UserUnfollowSfl) Make() UserUnfollowSflT {
	return func(username string, followedUsername string) (rpc.ProfileOut, error) {
		var zero rpc.ProfileOut
		user, _, err := s.UserFollowFl(username, followedUsername, false)
		if err != nil {
			return zero, err
		}
		profileOut := rpc.ProfileOut{}.FromModel(user, false)
		return profileOut, err
	}
}
