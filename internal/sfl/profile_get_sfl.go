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
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"

	"github.com/pvillela/go-foa-realworld/internal/arch/web"
)

// ProfileGetSflT is the type of the stereotype instance for the service flow that
// retrieves a user profile.
type ProfileGetSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	profileName string,
) (rpc.ProfileOut, error)

// ProfileGetSflC is the function that constructs a stereotype instance of type
// ProfileGetSflT without hard-wired stereotype dependencies.
func ProfileGetSflC(
	cfgSrc DefaultSflCfgSrc,
	userGetByNameDaf daf.UserGetByNameDafT,
	followingGetDaf daf.FollowingGetDafT,
) ProfileGetSflT {
	db := cfgSrc()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		profileName string,
	) (rpc.ProfileOut, error) {
		var zero rpc.ProfileOut
		var err error
		username := reqCtx.Username
		var currUser model.User
		var profileUser model.User
		follows := false

		if username != "" {
			currUser, err = userGetByNameDaf(ctx, tx, username)
			if err != nil {
				return zero, nil
			}
		}

		profileUser, err = userGetByNameDaf(ctx, tx, profileName)
		if err != nil {
			return zero, nil
		}

		following, err := followingGetDaf(ctx, tx, currUser.Id, profileUser.Id)
		if err != nil {
			return zero, nil
		}
		if following != (model.Following{}) {
			follows = true
		}

		profileOut := rpc.ProfileOut_FromModel(profileUser, follows)

		return profileOut, nil
	})
}
