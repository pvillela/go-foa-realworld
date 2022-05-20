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

// UserFollowSflT is the type of the stereotype instance for the service flow that
// causes the current user start following a given other user.
type UserFollowSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	followeeUsername string,
) (rpc.ProfileOut, error)

// UserFollowSflC is the function that constructs a stereotype instance of type
// UserFollowSflT.
func UserFollowSflC(
	ctxDb db.CtxDb,
) UserFollowSflT {
	followingCreateDaf := daf.FollowingCreateDafI
	return func(
		ctx context.Context,
		reqCtx web.RequestContext,
		followeeUsername string,
	) (rpc.ProfileOut, error) {
		username := reqCtx.Username
		var zero rpc.ProfileOut

		ctx, err := ctxDb.BeginTx(ctx)
		if err != nil {
			return zero, err
		}
		defer ctxDb.DeferredRollback(ctx)

		follower, _, err := daf.UserGetByNameDafI(ctx, username)
		if err != nil {
			return zero, err
		}

		followee, _, err := daf.UserGetByNameDafI(ctx, followeeUsername)
		if err != nil {
			return zero, err
		}

		tx, err := dbpgx.GetCtxTx(ctx)
		if err != nil {
			return zero, err
		}

		err = followingCreateDaf(ctx, tx, follower.Id, followee.Id)
		if err != nil {
			return zero, err
		}

		profileOut := rpc.ProfileOut_FromModel(&follower, true)
		return profileOut, nil
	}
}
