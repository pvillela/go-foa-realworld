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
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
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
		log.Error("transaction rollback failed ", err)
	}
}

func (s Db) WithTransaction(ctx context.Context, block func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := s.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer DeferredRollback(ctx, tx)

	err = block(ctx, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return errx.ErrxOf(err)
}

func DbWithTransaction[T any](
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
