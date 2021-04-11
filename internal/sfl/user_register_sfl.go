/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserRegisterSfl is the stereotype instance for the service flow that
// It represents the action of registering a user.
type UserRegisterSfl struct {
}

// UserRegisterSflT is the type of a function that takes an rpc.UserRegisterIn as input
// and returns a model.User.
type UserRegisterSflT = func(in rpc.UserRegisterIn) model.User
