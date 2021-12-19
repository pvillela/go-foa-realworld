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

// UserUnfollowSflT is the type of the stereotype instance for the service flow that
// causes the current user start following a given other user.
type UserUnfollowSflT = func(username string, followedUsername string) (rpc.ProfileOut, error)

// UserUnfollowSflC is the function that constructs a stereotype instance of type
// UserUnfollowSflT.
func UserUnfollowSflC(
	beginTxn func(context string) db.Txn,
	userFollowFl fs.UserFollowFlT,
) UserUnfollowSflT {
	return func(username string, followedUsername string) (rpc.ProfileOut, error) {
		txn := beginTxn("ArticleCreateSflS")
		defer txn.End()

		var zero rpc.ProfileOut
		user, _, err := userFollowFl(username, followedUsername, false, txn)
		if err != nil {
			return zero, err
		}
		profileOut := rpc.ProfileOut_FromModel(user, false)
		return profileOut, err
	}
}
