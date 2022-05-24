/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import "github.com/pvillela/go-foa-realworld/internal/model"

type ProfileOut struct {
	Profile model.Profile
}

func ProfileOut_FromModel(user *model.User, follows bool) ProfileOut {
	s := ProfileOut{}
	s.Profile = model.Profile_FromUser(user, follows)
	return s
}
