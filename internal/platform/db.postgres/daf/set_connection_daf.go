/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
)

type SetConnectionDafT = func(reqCtx context.Context)

func SetConnectionDafC(
	appCtx context.Context,
	connString string,
) SetConnectionDafT {
	pool, err := pgxpool.Connect(appCtx, connString)
	util.PanicOnError(err)
	return func(reqCtx context.Context) {
		conn, err := pool.Acquire(reqCtx)
		util.PanicOnError(err)
		SetCtxConn(reqCtx, conn)
	}
}
