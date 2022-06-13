/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package testutil

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	log "github.com/sirupsen/logrus"
	"testing"
)

var connStr = "postgres://testuser:testpassword@localhost:9999/testdb?sslmode=disable"

// DbTester runs txnlSubtest, providing a clean database.
func DbTester(t *testing.T, txnlSubtest dbpgx.TransactionalSubtest) {
	log.SetLevel(log.DebugLevel)

	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, connStr)
	errx.PanicOnError(err)

	db := dbpgx.Db{pool}
	defer pool.Close()

	CleanupAllTables(db, ctx)

	txnlSubtest(db, ctx, t)
}
