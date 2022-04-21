/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxCtxConnKeyT struct{}

var PgxCtxConnKey PgxCtxConnKeyT = struct{}{}

func SetCtxConn(ctx context.Context, conn *pgxpool.Conn) context.Context {
	return context.WithValue(ctx, PgxCtxConnKey, conn)
}

func GetCtxConn(ctx context.Context) *pgxpool.Conn {
	return ctx.Value(PgxCtxConnKey).(*pgxpool.Conn)
}

func ReleaseCtxConn(ctx context.Context) {
	conn := GetCtxConn(ctx)
	conn.Release()
}
