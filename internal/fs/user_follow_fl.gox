/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// UserFollowFlT is the type of the stereotype instance for the flow that
// causes the current user start following a given other user.
type UserFollowFlT = func(username string, followedUsername string, follow bool, txn db.Txn) (model.User, RecCtxUser, error)

// UserFollowFlC is the function that constructs a stereotype instance of type
// UserFollowFlT.
func UserFollowFlC(
	userGetByNameDaf UserGetByNameDafT,
	userUpdateDaf UserUpdateDafT,
) UserFollowFlT {
	return func(
		username string,
		followedUsername string,
		follow bool,
		txn db.Txn,
	) (model.User, RecCtxUser, error) {
		user, rc, err := userGetByNameDaf(username)
		if err != nil {
			return model.User{}, RecCtxUser{}, err
		}

		user = user.UpdateFollowees(followedUsername, follow)

		if rc, err = userUpdateDaf(user, rc, txn); err != nil {
			return model.User{}, RecCtxUser{}, err
		}

		return user, rc, nil
	}
}
