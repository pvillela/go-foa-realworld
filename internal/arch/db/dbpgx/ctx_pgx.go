/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package dbpgx

import (
	"context"
	"github.com/go-errors/errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	log "github.com/sirupsen/logrus"
)

type CtxPgxPoolKeyT struct{}
type CtxPgxConnKeyT struct{}
type CtxPgxTxKeyT struct{}

var CtxPgxPoolKey CtxPgxPoolKeyT = struct{}{}
var CtxPgxConnKey CtxPgxConnKeyT = struct{}{}
var CtxPgxTxKey CtxPgxTxKeyT = struct{}{}

type CtxPgx struct {
	Pool *pgxpool.Pool
}

func checkInterface(p CtxPgx) {
	var foo = func(cc db.CtxDb) {}
	foo(p)
}

func (p CtxPgx) SetPool(ctx context.Context) (context.Context, error) {
	if ctx.Value(CtxPgxPoolKey) != nil {
		return ctx, errors.New("ctx already has a Pool value")
	}
	return context.WithValue(ctx, CtxPgxPoolKey, p.Pool), nil
}

func GetCtxPool(ctx context.Context) (*pgxpool.Pool, error) {
	pool, ok := ctx.Value(CtxPgxPoolKey).(*pgxpool.Pool)
	var err error
	if !ok {
		err = errors.New("there is no Pool value in ctx")
	}
	return pool, err
}

func (CtxPgx) SetConn(ctx context.Context) (context.Context, error) {
	if ctx.Value(CtxPgxConnKey) != nil {
		return ctx, errors.New("ctx already has a connection value")
	}
	pool, err := GetCtxPool(ctx)
	if err != nil {
		return ctx, err
	}
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, CtxPgxConnKey, conn), nil
}

func GetCtxConn(ctx context.Context) (*pgxpool.Conn, error) {
	conn, ok := ctx.Value(CtxPgxConnKey).(*pgxpool.Conn)
	var err error
	if !ok {
		err = errors.New("there is no connection value in ctx")
	}
	return conn, err
}

func (CtxPgx) ReleaseConn(ctx context.Context) error {
	conn, err := GetCtxConn(ctx)
	if err != nil {
		return err
	}
	conn.Release()
	return nil
}

func (p CtxPgx) DeferredReleaseConn(ctx context.Context) {
	err := p.ReleaseConn(ctx)
	if err != nil {
		log.Error("connection release failed ", err)
	}
}

func (p CtxPgx) Begin(ctx context.Context) (context.Context, error) {
	if ctx.Value(CtxPgxConnKey) == nil {
		var err error
		ctx, err = p.SetConn(ctx)
		if err != nil {
			return nil, err
		}
	}
	conn, err := GetCtxConn(ctx)
	if err != nil {
		return ctx, err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, CtxPgxTxKey, tx), nil
}

func GetCtxTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(CtxPgxTxKey).(pgx.Tx)
	var err error
	if !ok {
		err = errors.New("there is no transaction value in ctx")
	}
	return tx, err
}

func (CtxPgx) Commit(ctx context.Context) error {
	tx, err := GetCtxTx(ctx)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (CtxPgx) Rollback(ctx context.Context) error {
	tx, err := GetCtxTx(ctx)
	if err != nil {
		return err
	}
	return tx.Rollback(ctx)
}

func (p CtxPgx) DeferredRollback(ctx context.Context) {
	err := p.Rollback(ctx)
	if err != nil {
		log.Error("rollback failed", err)
		p.DeferredReleaseConn(ctx)
	}
}
