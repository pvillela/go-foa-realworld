/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package dbpgxtest

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"strings"
	"testing"
)

// DafSubtest is the type of a function that implements a DAF subtest
// that is to be delimited by a transaction.
type DafSubtest func(
	ctx context.Context,
	tx pgx.Tx,
	t *testing.T,
)

// TransactionalSubtest is the tyype of a function that implements a DAF subtest
// that is already delimited by one or more transactions.
type TransactionalSubtest func(
	db dbpgx.Db,
	ctx context.Context,
	t *testing.T,
)

// TestPair pairs a TransactionalSubtest with a name for execution in a test suite.
type TestPair struct {
	Name string
	Func TransactionalSubtest
}

// TestWithTransaction is a convenience wrapper to transform a DafSubtest into a TransactionalSubtest.
func TestWithTransaction(
	f DafSubtest,
) TransactionalSubtest {
	return func(db dbpgx.Db, ctx context.Context, t *testing.T) {
		fL := util.Contextualize1V(dbpgx.WithTransaction[types.Unit], db, f)
		fL(ctx, t)
	}
}

// RunTestPairs executes a list of TestPair.
func RunTestPairs(db dbpgx.Db, ctx context.Context, t *testing.T, name string, testPairs []TestPair) {
	t.Run(name, func(t *testing.T) {
		for _, p := range testPairs {
			testFunc := func(t *testing.T) {
				defer util.Trace(p.Name)()
				p.Func(db, ctx, t)
			}
			t.Run(p.Name, testFunc)
		}
	})
}

// Parallel returns a decorated TransactionalSubtest that calls t.Parallel()
// just before executing txnlSubtest.
func Parallel(txnlSubtest TransactionalSubtest) TransactionalSubtest {
	return func(db dbpgx.Db, ctx context.Context, t *testing.T) {
		t.Parallel()
		txnlSubtest(db, ctx, t)
	}
}

// CleanupTables truncates tables. Should be called with all tables at once to reset
// the database while respecting referential integrity constraints.
func CleanupTables(ctx context.Context, tx pgx.Tx, tables ...string) {
	tablesStr := strings.Join(tables, ", ")
	sql := fmt.Sprintf("TRUNCATE %v", tablesStr)
	_, err := tx.Exec(ctx, sql)
	errx.PanicOnError(err)
}

// DbTester runs txnlSubtest with a clean database.
func DbTester(
	t *testing.T,
	txnlSubtest TransactionalSubtest,
	connStr string,
	cleanupTables func(db dbpgx.Db, ctx context.Context),
) {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, connStr)
	errx.PanicOnError(err)

	db := dbpgx.Db{pool}
	defer pool.Close()

	cleanupTables(db, ctx)

	txnlSubtest(db, ctx, t)
}
