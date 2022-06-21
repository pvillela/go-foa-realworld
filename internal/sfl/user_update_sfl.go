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
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserUpdateSflT is the type of the stereotype instance for the service flow that
// updates a user.
type UserUpdateSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	in rpc.UserUpdateIn,
) (rpc.UserOut, error)

// UserUpdateSflC is the function that constructs a stereotype instance of type
// UserUpdateSflT with hard-wired stereotype dependencies.
func UserUpdateSflC(
	ctxDb cdb.CtxDb,
) UserUpdateSflT {
	userGetByNameDaf := daf.UserGetByNameDaf
	userUpdateDaf := daf.UserUpdateDaf
	return UserUpdateSflC0(
		ctxDb,
		userGetByNameDaf,
		userUpdateDaf,
	)
}

// UserUpdateSflC0 is the function that constructs a stereotype instance of type
// UserUpdateSflT without hard-wired stereotype dependencies.
func UserUpdateSflC0(
	ctxDb cdb.CtxDb,
	userGetByNameDaf daf.UserGetByNameDafT,
	userUpdateDaf daf.UserUpdateDafT,
) UserUpdateSflT {
	return cdb.SflWithTransaction(ctxDb, func(
		ctx context.Context,
		reqCtx web.RequestContext,
		in rpc.UserUpdateIn,
	) (rpc.UserOut, error) {
		username := reqCtx.Username

		user, rc, err := userGetByNameDaf(ctx, username)
		if err != nil {
			return rpc.UserOut{}, err
		}

		user = user.Update(in.User)

		_, err = userUpdateDaf(ctx, user, rc)
		if err != nil {
			return rpc.UserOut{}, err
		}

		token := reqCtx.Token

		userOut := rpc.UserOut_FromModel(user, token.Raw)
		return userOut, nil
	})
}
