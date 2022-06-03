/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"

	"github.com/pvillela/go-foa-realworld/internal/arch/db/cdb"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

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
// UserFollowSflT with hard-wired stereotype dependencies.
func UserFollowSflC(
	ctxDb cdb.CtxDb,
) UserFollowSflT {
	userGetByNameDaf := daf.UserGetByNameDafI
	followingCreateDaf := daf.FollowingCreateDafI
	return UserFollowSflC0(
		ctxDb,
		userGetByNameDaf,
		followingCreateDaf,
	)
}

// UserFollowSflC0 is the function that constructs a stereotype instance of type
// UserFollowSflT without hard-wired stereotype dependencies.
func UserFollowSflC0(
	ctxDb cdb.CtxDb,
	userGetByNameDaf daf.UserGetByNameDafT,
	followingCreateDaf daf.FollowingCreateDafT,
) UserFollowSflT {
	return func(
		ctx context.Context,
		reqCtx web.RequestContext,
		followeeUsername string,
	) (rpc.ProfileOut, error) {
		username := reqCtx.Username

		return cdb.WithTransaction(ctxDb, ctx, func(
			ctx context.Context,
		) (rpc.ProfileOut, error) {
			var zero rpc.ProfileOut

			follower, _, err := userGetByNameDaf(ctx, username)
			if err != nil {
				return zero, err
			}

			followee, _, err := userGetByNameDaf(ctx, followeeUsername)
			if err != nil {
				return zero, err
			}

			tx, err := dbpgx.GetCtxTx(ctx)
			if err != nil {
				return zero, err
			}

			_, err = followingCreateDaf(ctx, tx, follower.Id, followee.Id)
			if err != nil {
				return zero, err
			}

			profileOut := rpc.ProfileOut_FromModel(follower, true)
			return profileOut, nil
		})
	}
}
