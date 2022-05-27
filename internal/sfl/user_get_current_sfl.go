/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserGetCurrentSflT is the type of the stereotype instance for the service flow that
// returns the current user.
type UserGetCurrentSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	_ types.Unit,
) (rpc.UserOut, error)

// UserGetCurrentSflC is the function that constructs a stereotype instance of type
// UserGetCurrentSflT with hard-wired stereotype dependencies.
func UserGetCurrentSflC(
	ctxDb db.CtxDb,
) UserGetCurrentSflT {
	userGetByNameDaf := daf.UserGetByNameDafI
	return UserGetCurrentSflC0(
		ctxDb,
		userGetByNameDaf,
	)
}

// UserGetCurrentSflC0 is the function that constructs a stereotype instance of type
// UserGetCurrentSflT without hard-wired stereotype dependencies.
func UserGetCurrentSflC0(
	ctxDb db.CtxDb,
	userGetByNameDaf daf.UserGetByNameDafT,
) UserGetCurrentSflT {
	return func(
		ctx context.Context,
		reqCtx web.RequestContext,
		_ types.Unit,
	) (rpc.UserOut, error) {
		return db.WithTransaction(ctxDb, ctx, func(
			ctx context.Context,
		) (rpc.UserOut, error) {
			username := reqCtx.Username

			user, _, err := userGetByNameDaf(ctx, username)
			if err != nil {
				return rpc.UserOut{}, err
			}

			token := reqCtx.Token

			userOut := rpc.UserOut_FromModel(user, token.Raw)
			return userOut, nil
		})
	}
}
