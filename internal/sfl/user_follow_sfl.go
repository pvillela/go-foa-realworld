/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentAddSflS is the stereotype instance for the service flow that
// causes the current user start following a given other user.
type UserFollowSflS struct {
	BeginTxn     func(context string) db.Txn
	UserFollowFl fs.UserFollowFlT
}

// UserFollowSflT is the function type instantiated by UserFollowSflS.
type UserFollowSflT = func(username string, followedUsername string) (rpc.ProfileOut, error)

func (s UserFollowSflS) Make() UserFollowSflT {
	return func(username string, followedUsername string) (rpc.ProfileOut, error) {
		txn := s.BeginTxn("ArticleCreateSflS")
		defer txn.End()

		var zero rpc.ProfileOut
		user, _, err := s.UserFollowFl(username, followedUsername, true, txn)
		if err != nil {
			return zero, err
		}
		profileOut := rpc.ProfileOut_FromModel(user, true)
		return profileOut, err
	}
}
