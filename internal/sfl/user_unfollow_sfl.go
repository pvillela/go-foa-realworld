/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserUnfollowSflT is the type of the stereotype instance for the service flow that
// causes the current user to stop following a given other user.
type UserUnfollowSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	followeeUsername string,
) (rpc.ProfileOut, error)

// UserUnfollowSflC is the function that constructs a stereotype instance of type
// UserUnfollowSflT with hard-wired stereotype dependencies.
func UserUnfollowSflC(
	ctxDb db.CtxDb,
) UserFollowSflT {
	userGetByNameDaf := daf.UserGetByNameDafI
	followingDeleteDaf := daf.FollowingDeleteDafI
	return UserUnfollowSflC0(
		ctxDb,
		userGetByNameDaf,
		followingDeleteDaf,
	)
}

// UserUnfollowSflC0 is the function that constructs a stereotype instance of type
// UserUnfollowSflT without hard-wired stereotype dependencies.
func UserUnfollowSflC0(
	ctxDb db.CtxDb,
	userGetByNameDaf daf.UserGetByNameDafT,
	followingDeleteDaf daf.FollowingDeleteDafT,
) UserFollowSflT {
	return func(
		ctx context.Context,
		reqCtx web.RequestContext,
		followeeUsername string,
	) (rpc.ProfileOut, error) {
		username := reqCtx.Username

		return db.WithTransaction(ctxDb, ctx, func(
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

			err = followingDeleteDaf(ctx, tx, follower.Id, followee.Id)
			if err != nil {
				return zero, err
			}

			profileOut := rpc.ProfileOut_FromModel(&follower, true)
			return profileOut, nil
		})
	}
}
