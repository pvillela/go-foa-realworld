/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserUpdateSfl is the stereotype instance for the service flow that
// It represents the action of registering a user.
type UserUpdateSfl struct {
	UserGetByNameDaf fs.UserGetByNameDafT
	UserUpdateDaf    fs.UserUpdateDafT
	UserGenTokenBf   fs.UserGenTokenBfT
}

// UserUpdateSflT is the type of a function that takes an rpc.UserUpdateIn as input
// and returns a model.User.
type UserUpdateSflT = func(username string, in rpc.UserUpdateIn) (rpc.UserOut, error)

func (s UserUpdateSfl) Make() UserUpdateSflT {
	return func(username string, in rpc.UserUpdateIn) (rpc.UserOut, error) {
		user, rc, err := s.UserGetByNameDaf(username)
		if err != nil {
			return rpc.UserOut{}, err
		}

		fieldsToUpdate := make(map[model.UserUpdatableField]interface{}, 5)
		if v := in.User.Username; v != nil {
			fieldsToUpdate[model.UserName] = *v
		}
		if v := in.User.Email; v != nil {
			fieldsToUpdate[model.UserEmail] = *v
		}
		if v := in.User.Password; v != nil {
			fieldsToUpdate[model.UserPassword] = *v
		}
		if v := in.User.Bio; v != nil {
			fieldsToUpdate[model.UserBio] = v
		}
		if v := in.User.Image; v != nil {
			fieldsToUpdate[model.UserImageLink] = *v
		}

		user = user.Update(fieldsToUpdate)

		_, err = s.UserUpdateDaf(user, rc)
		if err != nil {
			return rpc.UserOut{}, err
		}

		token, err := s.UserGenTokenBf(user)
		if err != nil {
			return rpc.UserOut{}, err
		}

		userOut := rpc.UserOut{}.FromModel(user, token)
		return userOut, err
	}
}
