/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UnfollowUserSflS contains the dependencies required for the construction of a
// UnfollowUserSfl. It represents the action of having the current user stop following a given
// other user.
type UnfollowUserSflS struct {
}

// UnfollowUserSfl is the type of a function that takes the current username and a followed
// username and returns a model.ProfileOut.
type UnfollowUserSfl = func(currentUsername string, followedUsername string) rpc.ProfileOut
