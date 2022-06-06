/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/cdb"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserRegisterSflT is the type of the stereotype instance for the service flow that
// represents the action of registering a user.
type UserRegisterSflT = func(
	ctx context.Context,
	_ web.RequestContext,
	in rpc.UserRegisterIn,
) (rpc.UserOut, error)

// UserRegisterSflC is the function that constructs a stereotype instance of type
// UserRegisterSflT with hard-wired stereotype dependencies.
func UserRegisterSflC(
	ctxDb cdb.CtxDb,
	userGenTokenBf bf.UserGenTokenBfT,
) UserRegisterSflT {
	userCreateDaf := daf.UserCreateDafI
	return UserRegisterSflC0(
		ctxDb,
		userCreateDaf,
		userGenTokenBf,
	)
}

// UserRegisterSflC0 is the function that constructs a stereotype instance of type
// UserRegisterSflT without hard-wired stereotype dependencies.
func UserRegisterSflC0(
	ctxDb cdb.CtxDb,
	userCreateDaf daf.UserCreateDafT,
	userGenTokenBf bf.UserGenTokenBfT,
) UserRegisterSflT {
	return cdb.SflWithTransaction(ctxDb, func(
		ctx context.Context,
		reqCtx web.RequestContext,
		in rpc.UserRegisterIn,
	) (rpc.UserOut, error) {
		user := in.ToUser()

		_, err := userCreateDaf(ctx, &user)
		if err != nil {
			return rpc.UserOut{}, err
		}

		token, err := userGenTokenBf(user)
		if err != nil {
			return rpc.UserOut{}, err
		}

		userOut := rpc.UserOut_FromModel(user, token)
		return userOut, nil
	})
}
