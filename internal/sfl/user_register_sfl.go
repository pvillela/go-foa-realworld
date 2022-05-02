/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserRegisterSflT is the type of the stereotype instance for the service flow that
// represents the action of registering a user, with hard-wired BF dependencies.
type UserRegisterSflT = func(_ context.Context, in rpc.UserRegisterIn) (rpc.UserOut, error)

// UserRegisterSflC is the function that constructs a stereotype instance of type
// UserRegisterSflT with hard-wired BF dependencies.
func UserRegisterSflC(
	beginTxn func(context string) db.Txn,
	userCreateDaf daf.UserCreateDafT,
) UserRegisterSflT {
	userGenTokenBf := fs.UserGenTokenBfI
	return UserRegisterSflC0(
		beginTxn,
		userCreateDaf,
		userGenTokenBf,
	)
}

// UserRegisterSflC0 is the function that constructs a stereotype instance of type
// UserRegisterSflT without hard-wired BF dependencies.
func UserRegisterSflC0(
	beginTxn func(context string) db.Txn,
	userCreateDaf daf.UserCreateDafT,
	userGenTokenBf fs.UserGenTokenBfT,
) UserRegisterSflT {
	return func(_ context.Context, in rpc.UserRegisterIn) (rpc.UserOut, error) {
		txn := beginTxn("UserRegisterSflT")
		defer txn.End()

		user := in.ToUser()

		_, err := userCreateDaf(user, txn)
		if err != nil {
			return rpc.UserOut{}, err
		}

		token, err := userGenTokenBf(user)
		if err != nil {
			return rpc.UserOut{}, err
		}

		userOut := rpc.UserOut_FromModel(user, token)
		return userOut, err
	}
}
