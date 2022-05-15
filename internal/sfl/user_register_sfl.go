/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserRegisterSflT is the type of the stereotype instance for the service flow that
// represents the action of registering a user.
type UserRegisterSflT = func(ctx context.Context, in rpc.UserRegisterIn) (rpc.UserOut, error)

// UserRegisterSflC is the function that constructs a stereotype instance of type
// UserRegisterSflT.
func UserRegisterSflC(
	ctxDb db.CtxDb,
	userCreateDaf daf.UserCreateDafT,
	userGenTokenBf bf.UserGenTokenBfT,
) UserRegisterSflT {
	return func(ctx context.Context, in rpc.UserRegisterIn) (rpc.UserOut, error) {
		ctx, err := ctxDb.BeginTx(ctx)
		if err != nil {
			return rpc.UserOut{}, err
		}
		defer ctxDb.DeferredRollback(ctx)

		user := in.ToUser()

		_, err = userCreateDaf(ctx, user)
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
