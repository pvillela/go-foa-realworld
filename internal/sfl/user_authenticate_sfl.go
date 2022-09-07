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
	"github.com/pvillela/go-foa-realworld/internal/daf"
	rpc2 "github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserAuthenticateSflT is the type of the stereotype instance for the service flow that
// authenticates a user.
type UserAuthenticateSflT = func(
	ctx context.Context,
	_ web.RequestContext,
	in rpc2.UserAuthenticateIn,
) (rpc2.UserOut, error)

// UserAuthenticateSflC is the function that constructs a stereotype instance of type
// UserAuthenticateSflT with hard-wired stereotype dependencies.
func UserAuthenticateSflC(
	cfgSrc UserSflCfgSrc,
	userGenTokenBf bf.UserGenTokenBfT,
) UserAuthenticateSflT {
	return UserAuthenticateSflC0(
		cfgSrc,
		daf.UserGetByEmailDaf,
		userGenTokenBf,
		bf.UserAuthenticateBf,
	)
}

// UserAuthenticateSflC0 is the function that constructs a stereotype instance of type
// UserAuthenticateSflT without hard-wired stereotype dependencies.
func UserAuthenticateSflC0(
	cfgSrc UserSflCfgSrc,
	userGetByEmailDaf daf.UserGetByEmailDafT,
	userGenTokenBf bf.UserGenTokenBfT,
	userAuthenticateBf bf.UserAuthenticateBfT,
) UserAuthenticateSflT {
	ctxDb := cfgSrc()
	return cdb.SflWithTransaction(ctxDb, func(
		ctx context.Context,
		reqCtx web.RequestContext,
		in rpc2.UserAuthenticateIn,
	) (rpc2.UserOut, error) {
		var zero rpc2.UserOut
		email := in.User.Email
		password := in.User.Password

		user, _, err := userGetByEmailDaf(ctx, email)
		if err != nil {
			return zero, err
		}

		if !userAuthenticateBf(user, password) {
			// The error info below is not secure but OK for now for debugging
			return zero, bf.ErrAuthenticationFailed.Make(nil,
				bf.ErrMsgAuthenticationFailed, user.Username, password)
		}

		token, err := userGenTokenBf(user)
		if err != nil {
			return zero, err
		}

		userOut := rpc2.UserOut_FromModel(user, token)
		return userOut, err
	})
}
