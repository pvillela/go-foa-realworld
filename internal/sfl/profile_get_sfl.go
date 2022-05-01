/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/newdaf"

	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ProfileGetSflT is the type of the stereotype instance for the service flow that
// retrieves a user profile.
type ProfileGetSflT = func(ctx context.Context, profileName string) (rpc.ProfileOut, error)

// ProfileGetSflC is the function that constructs a stereotype instance of type
// ProfileGetSflT.
func ProfileGetSflC(
	userGetByNameDaf newdaf.UserGetByNameDafT,
) ProfileGetSflT {
	return func(ctx context.Context, profileName string) (rpc.ProfileOut, error) {
		username := web.ContextToRequestContext(ctx).Username

		var zero rpc.ProfileOut

		profileUser, _, err := userGetByNameDaf(profileName)
		if err != nil {
			return zero, err
		}

		var follows bool
		if username != "" {
			user, _, err := userGetByNameDaf(username)
			if err != nil {
				return zero, err
			}
			follows = user.Follows(profileName)
		}

		profileOut := rpc.ProfileOut_FromModel(profileUser, follows)

		return profileOut, nil
	}
}
