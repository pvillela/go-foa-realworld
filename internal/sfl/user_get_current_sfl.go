/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"

	"github.com/pvillela/go-foa-realworld/internal/arch"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleCreateSflT is the type of the stereotype instance for the service flow that
// returns the current user.
type UserGetCurrentSflT = func(ctx context.Context, _ arch.Unit) (rpc.UserOut, error)

// UserGetCurrentSflC is the function that constructs a stereotype instance of type
// UserGetCurrentSflT.
func UserGetCurrentSflC(
	userGetByNameDaf daf.UserGetByNameDafT,
) UserGetCurrentSflT {
	return func(ctx context.Context, _ arch.Unit) (rpc.UserOut, error) {
		username := web.ContextToRequestContext(ctx).Username

		user, _, err := userGetByNameDaf(username)
		if err != nil {
			return rpc.UserOut{}, err
		}

		token := web.ContextToRequestContext(ctx).Token

		userOut := rpc.UserOut_FromModel(user, token.Raw)
		return userOut, err
	}
}
