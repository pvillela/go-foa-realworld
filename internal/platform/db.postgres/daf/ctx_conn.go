/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxCtxConnKeyT struct{}
type PgxCtxTxKeyT struct{}

var PgxCtxConnKey PgxCtxConnKeyT = struct{}{}
var PgxCtxTxKey PgxCtxTxKeyT = struct{}{}

func SetCtxConn(ctx context.Context, pool *pgxpool.Pool) (context.Context, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, PgxCtxConnKey, conn), nil
}

func GetCtxConn(ctx context.Context) (*pgxpool.Conn, error) {
	conn, ok := ctx.Value(PgxCtxConnKey).(*pgxpool.Conn)
	var err error
	if !ok {
		err = errors.New("there is no connection value in ctx")
	}
	return conn, err
}

func ReleaseCtxConn(ctx context.Context) error {
	conn, err := GetCtxConn(ctx)
	if err != nil {
		return err
	}
	conn.Release()
	return nil
}

func BeginTx(ctx context.Context) (context.Context, error) {
	conn, err := GetCtxConn(ctx)
	if err != nil {
		return ctx, err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, PgxCtxTxKey, tx), nil
}

func BeginConnTx(ctx context.Context, pool *pgxpool.Pool) (context.Context, error) {
	ctx, err := SetCtxConn(ctx, pool)
	if err != nil {
		return ctx, err
	}
	return BeginTx(ctx)
}

func GetCtxTx(ctx context.Context) (*pgxpool.Tx, error) {
	tx, ok := ctx.Value(PgxCtxTxKey).(*pgxpool.Tx)
	var err error
	if !ok {
		err = errors.New("there is no transaction value in ctx")
	}
	return tx, err
}

func CommitTx(ctx context.Context) error {
	tx, err := GetCtxTx(ctx)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func RollbackTx(ctx context.Context) error {
	tx, err := GetCtxTx(ctx)
	if err != nil {
		return err
	}
	return tx.Rollback(ctx)
}
