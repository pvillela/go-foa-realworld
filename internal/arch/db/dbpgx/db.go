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