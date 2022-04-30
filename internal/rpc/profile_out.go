/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import "github.com/pvillela/go-foa-realworld/internal/model"

type Profile struct {
	Username  string  `json:"username"`
	Bio       *string `json:"bio"`
	Image     *string `json:"image"`
	Following bool    `json:"following"`
}

type ProfileOut struct {
	Profile Profile
}

func Profile_FromModel(user *model.User, follows bool) Profile {
	s := Profile{}
	if user.Bio != nil {
		s.Bio = user.Bio
	} else {
		empty := ""
		s.Bio = &empty
	}

	if user.ImageLink != "" {
		s.Image = &user.ImageLink
	}

	s.Username = user.Username
	s.Following = follows

	return s
}

func ProfileOut_FromModel(user *model.User, follows bool) ProfileOut {
	s := ProfileOut{}
	s.Profile = Profile_FromModel(user, follows)
	return s
}
