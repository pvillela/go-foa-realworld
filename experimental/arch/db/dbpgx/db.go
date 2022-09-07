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
	"github.com/pvillela/go-foa-realworld/experimental/arch/errx"
	"github.com/pvillela/go-foa-realworld/experimental/arch/util"
	log "github.com/sirupsen/logrus"
)

type Db struct {
	Pool *pgxpool.Pool
}

func (s Db) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	conn, err := s.Pool.Acquire(ctx)
	return conn, errx.ErrxOf(err)
}

func (s Db) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: "Serializable"})
	return tx, errx.ErrxOf(err)
}

func DeferredRollback(ctx context.Context, tx pgx.Tx) {
	err := tx.Rollback(ctx)
	if err != nil {
		log.Debug("transaction rollback failed ", err)
	}
}

func WithTransaction[T any](
	db Db,
	ctx context.Context,
	block func(ctx context.Context, tx pgx.Tx) (T, error),
) (T, error) {
	var zero T
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return zero, err
	}

	defer DeferredRollback(ctx, tx)

	t, err := block(ctx, tx)
	if err != nil {
		return zero, err
	}

	err = tx.Commit(ctx)
	return t, errx.ErrxOf(err)
}

func SflWithTransaction[R, S, T any](
	db Db,
	block func(ctx context.Context, tx pgx.Tx, reqCtx R, in S) (T, error),
) func(ctx context.Context, reqCtx R, in S) (T, error) {
	return util.Contextualize2(WithTransaction[T], db, block)
}

// Implementation from scratch of above, without use of util.Contextualize2
func sflWithTransaction0[R, S, T any](
	db Db,
	block func(ctx context.Context, tx pgx.Tx, reqCtx R, in S) (T, error),
) func(ctx context.Context, reqCtx R, in S) (T, error) {
	return func(ctx context.Context, reqCtx R, in S) (T, error) {
		block1 := func(ctx context.Context, tx pgx.Tx) (T, error) {
			return block(ctx, tx, reqCtx, in)
		}
		return WithTransaction(db, ctx, block1)
	}
}
