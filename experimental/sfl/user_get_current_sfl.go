/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/experimental/arch/db/cdb"
	"github.com/pvillela/go-foa-realworld/experimental/arch/types"
	"github.com/pvillela/go-foa-realworld/experimental/arch/web"
	"github.com/pvillela/go-foa-realworld/experimental/daf"
	"github.com/pvillela/go-foa-realworld/experimental/rpc"
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
	cfgSrc UserSflCfgSrc,
) UserGetCurrentSflT {
	return UserGetCurrentSflC0(
		cfgSrc,
		daf.UserGetByNameDaf,
	)
}

// UserGetCurrentSflC0 is the function that constructs a stereotype instance of type
// UserGetCurrentSflT without hard-wired stereotype dependencies.
func UserGetCurrentSflC0(
	cfgSrc UserSflCfgSrc,
	userGetByNameDaf daf.UserGetByNameDafT,
) UserGetCurrentSflT {
	ctxDb := cfgSrc()
	return cdb.SflWithTransaction(ctxDb, func(
		ctx context.Context,
		reqCtx web.RequestContext,
		_ types.Unit,
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
