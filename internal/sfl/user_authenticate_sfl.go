/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserAuthenticateSflT is the type of the stereotype instance for the service flow that
// authenticates a user.
type UserAuthenticateSflT = func(
	ctx context.Context,
	_ web.RequestContext,
	in rpc.UserAuthenticateIn,
) (rpc.UserOut, error)

// UserAuthenticateSflC is the function that constructs a stereotype instance of type
// UserAuthenticateSflT with hard-wired stereotype dependencies.
func UserAuthenticateSflC(
	ctxDb db.CtxDb,
	userGenTokenBf bf.UserGenTokenBfT,
) UserAuthenticateSflT {
	userGetByEmailDaf := daf.UserGetByEmailDafI
	userAuthenticateBf := bf.UserAuthenticateBfI
	return UserAuthenticateSflC0(
		ctxDb,
		userGetByEmailDaf,
		userGenTokenBf,
		userAuthenticateBf,
	)
}

// UserAuthenticateSflC0 is the function that constructs a stereotype instance of type
// UserAuthenticateSflT without hard-wired stereotype dependencies.
func UserAuthenticateSflC0(
	ctxDb db.CtxDb,
	userGetByEmailDaf daf.UserGetByEmailDafT,
	userGenTokenBf bf.UserGenTokenBfT,
	userAuthenticateBf bf.UserAuthenticateBfT,
) UserAuthenticateSflT {
	return func(
		ctx context.Context,
		_ web.RequestContext,
		in rpc.UserAuthenticateIn,
	) (rpc.UserOut, error) {
		block := func(ctx context.Context) (rpc.UserOut, error) {
			var zero rpc.UserOut
			email := in.User.Email
			password := in.User.Password

			user, _, err := userGetByEmailDaf(ctx, email)
			if err != nil {
				return zero, err
			}

			if !userAuthenticateBf(user, password, user.PasswordSalt) {
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
		}

		return db.CtxDb_WithTransaction(ctxDb, ctx, block)
	}
}
