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

// UserFollowSflT is the type of the stereotype instance for the service flow that
// causes the current user start following a given other user.
type UserFollowSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	followeeUsername string,
) (rpc.ProfileOut, error)

// UserFollowSflC is the function that constructs a stereotype instance of type
// UserFollowSflT with configuration information and hard-wired stereotype dependencies.
func UserFollowSflC(
	cfgSrc DefaultSflCfgSrc,
) UserFollowSflT {
	return UserFollowSflC0(
		cfgSrc,
		daf.UserGetByNameDaf,
		daf.FollowingCreateDaf,
	)
}

// UserFollowSflC0 is the function that constructs a stereotype instance of type
// UserFollowSflT without hard-wired stereotype dependencies.
func UserFollowSflC0(
	cfgSrc DefaultSflCfgSrc,
	userGetByNameDaf daf.UserGetByNameDafT,
	followingCreateDaf daf.FollowingCreateDafT,
) UserFollowSflT {
	db := cfgSrc()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		followeeUsername string,
	) (rpc.ProfileOut, error) {
		username := reqCtx.Username
		var zero rpc.ProfileOut

		follower, err := userGetByNameDaf(ctx, tx, username)
		if err != nil {
			return zero, err
		}

		followee, err := userGetByNameDaf(ctx, tx, followeeUsername)
		if err != nil {
			return zero, err
		}

		_, err = followingCreateDaf(ctx, tx, follower.Id, followee.Id)
		if err != nil {
			return zero, err
		}

		profileOut := rpc.ProfileOut_FromModel(followee, true)
		return profileOut, nil
	})
}
