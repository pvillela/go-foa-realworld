/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package dbpgx

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pvillela/go-foa-realworld/experimental/arch/db/cdb"
	"github.com/pvillela/go-foa-realworld/experimental/arch/errx"
	log "github.com/sirupsen/logrus"
)

type CtxPgxPoolKeyT struct{}
type CtxPgxTxKeyT struct{}

var CtxPgxPoolKey CtxPgxPoolKeyT = struct{}{}
var CtxPgxTxKey CtxPgxTxKeyT = struct{}{}

type CtxPgx struct {
	Pool *pgxpool.Pool
}

// Interface verification
func _(p CtxPgx) {
	func(cdb cdb.CtxDb) {}(p)
}

func (p CtxPgx) SetPool(ctx context.Context) (context.Context, error) {
	if ctx.Value(CtxPgxPoolKey) != nil {
		return ctx, errx.NewErrx(nil, "ctx already has a Pool value")
	}
	return context.WithValue(ctx, CtxPgxPoolKey, p.Pool), nil
}

func GetCtxPool(ctx context.Context) (*pgxpool.Pool, error) {
	pool, ok := ctx.Value(CtxPgxPoolKey).(*pgxpool.Pool)
	var err error
	if !ok {
		err = errx.NewErrx(nil, "there is no Pool value in ctx")
	}
	return pool, err
}

func (p CtxPgx) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := p.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: "Serializable"})
	if err != nil {
		return ctx, errx.ErrxOf(err)
	}
	return context.WithValue(ctx, CtxPgxTxKey, tx), nil
}

func GetCtxTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(CtxPgxTxKey).(pgx.Tx)
	var err error
	if !ok {
		err = errx.NewErrx(nil, "there is no transaction value in ctx")
	}
	return tx, err
}

func (p CtxPgx) Commit(ctx context.Context) (context.Context, error) {
	tx, err := GetCtxTx(ctx)
	if err != nil {
		return ctx, errx.ErrxOf(err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return ctx, errx.ErrxOf(err)
	}
	ctx = context.WithValue(ctx, CtxPgxTxKey, nil)
	return ctx, nil
}

func (p CtxPgx) Rollback(ctx context.Context) (context.Context, error) {
	tx, err := GetCtxTx(ctx)
	if err != nil {
		return ctx, errx.ErrxOf(err)
	}
	err = tx.Rollback(ctx)
	if err != nil {
		return ctx, errx.ErrxOf(err)
	}

	ctx = context.WithValue(ctx, CtxPgxTxKey, nil)
	return ctx, nil
}

func (p CtxPgx) DeferredRollback(ctx context.Context) {
	_, err := p.Rollback(ctx)
	if err != nil {
		log.Debug("rollback failed ", err)
	}
}
