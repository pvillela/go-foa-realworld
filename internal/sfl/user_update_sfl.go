/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// UserUpdateSflT is the type of the stereotype instance for the service flow that
// updates a user.
type UserUpdateSflT = func(ctx context.Context, in rpc.UserUpdateIn) (rpc.UserOut, error)

// UserUpdateSflC is the function that constructs a stereotype instance of type
// UserUpdateSflT.
func UserUpdateSflC(
	ctxDb db.CtxDb,
	userGetByNameDaf fs.UserGetByNameDafT,
	userUpdateDaf fs.UserUpdateDafT,
) UserUpdateSflT {
	return func(ctx context.Context, in rpc.UserUpdateIn) (rpc.UserOut, error) {
		username := web.ContextToRequestContext(ctx).Username

		ctx, err := ctxDb.BeginTx(ctx)
		if err != nil {
			return rpc.UserOut{}, err
		}
		defer ctxDb.DeferredRollback(ctx)

		user, rc, err := userGetByNameDaf(ctx, username)
		if err != nil {
			return rpc.UserOut{}, err
		}

		user = user.Update(in.User)

		_, err = userUpdateDaf(ctx, user, rc)
		if err != nil {
			return rpc.UserOut{}, err
		}

		_, err = ctxDb.Commit(ctx)
		if err != nil {
			return rpc.UserOut{}, err
		}

		token := web.ContextToRequestContext(ctx).Token

		userOut := rpc.UserOut_FromModel(user, token.Raw)
		return userOut, nil
	}
}
