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
	"github.com/pvillela/go-foa-realworld/internal/daf"
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
// UserUpdateSflT with configuration information and hard-wired stereotype dependencies.
func UserUpdateSflC(
	cfgSrc DefaultSflCfgSrc,
) UserUpdateSflT {
	return UserUpdateSflC0(
		cfgSrc,
		daf.UserGetByNameDaf,
		daf.UserUpdateDaf,
	)
}

// UserUpdateSflC0 is the function that constructs a stereotype instance of type
// UserUpdateSflT without hard-wired stereotype dependencies.
func UserUpdateSflC0(
	cfgSrc DefaultSflCfgSrc,
	userGetByNameDaf daf.UserGetByNameDafT,
	userUpdateDaf daf.UserUpdateDafT,
) UserUpdateSflT {
	db := cfgSrc()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		in rpc.UserUpdateIn,
	) (rpc.UserOut, error) {
		username := reqCtx.Username

		user, err := userGetByNameDaf(ctx, tx, username)
		if err != nil {
			return rpc.UserOut{}, err
		}

		user = user.Update(in.User)

		err = userUpdateDaf(ctx, tx, &user)
		if err != nil {
			return rpc.UserOut{}, err
		}

		token := reqCtx.Token

		userOut := rpc.UserOut_FromModel(user, token.Raw)
		return userOut, nil
	})
}
