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
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	rpc2 "github.com/pvillela/go-foa-realworld/rpc"
)

// UserUpdateSflT is the type of the stereotype instance for the service flow that
// updates a user.
type UserUpdateSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	in rpc2.UserUpdateIn,
) (rpc2.UserOut, error)

// UserUpdateSflC is the function that constructs a stereotype instance of type
// UserUpdateSflT with hard-wired stereotype dependencies.
func UserUpdateSflC(
	cfgPvdr UserSflCfgPvdr,
) UserUpdateSflT {
	return UserUpdateSflC0(
		cfgPvdr,
		daf.UserGetByNameDaf,
		daf.UserUpdateDaf,
	)
}

// UserUpdateSflC0 is the function that constructs a stereotype instance of type
// UserUpdateSflT without hard-wired stereotype dependencies.
func UserUpdateSflC0(
	cfgPvdr UserSflCfgPvdr,
	userGetByNameDaf daf.UserGetByNameDafT,
	userUpdateDaf daf.UserUpdateDafT,
) UserUpdateSflT {
	ctxDb := cfgPvdr()
	return cdb.SflWithTransaction(ctxDb, func(
		ctx context.Context,
		reqCtx web.RequestContext,
		in rpc2.UserUpdateIn,
	) (rpc2.UserOut, error) {
		username := reqCtx.Username

		user, rc, err := userGetByNameDaf(ctx, username)
		if err != nil {
			return rpc2.UserOut{}, err
		}

		user = user.Update(in.User)

		_, err = userUpdateDaf(ctx, user, rc)
		if err != nil {
			return rpc2.UserOut{}, err
		}

		token := reqCtx.Token

		userOut := rpc2.UserOut_FromModel(user, token.Raw)
		return userOut, nil
	})
}
