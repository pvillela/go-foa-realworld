/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserAuthenticateSflT is the type of the stereotype instance for the service flow that
// authenticates a user.
type UserAuthenticateSflT = func(
	ctx context.Context,
	_ web.RequestContext,
	in rpc.UserAuthenticateIn,
) (rpc.UserOut, error)

// UserAuthenticateSflC0 is the function that constructs a stereotype instance of type
// UserAuthenticateSflT without hard-wired stereotype dependencies.
func UserAuthenticateSflC0(
	cfgSrc DefaultSflCfgSrc,
	userGetByEmailDaf daf.UserGetByEmailDafT,
	userGenTokenBf bf.UserGenTokenBfT,
	userAuthenticateBf bf.UserAuthenticateBfT,
) UserAuthenticateSflT {
	db := cfgSrc()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		in rpc.UserAuthenticateIn,
	) (rpc.UserOut, error) {
		var zero rpc.UserOut
		email := in.User.Email
		password := in.User.Password

		user, err := userGetByEmailDaf(ctx, tx, email)
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

		userOut := rpc.UserOut_FromModel(user, token)
		return userOut, err
	})
}
