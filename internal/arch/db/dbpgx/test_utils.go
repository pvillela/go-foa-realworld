/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package dbpgx

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
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
	db Db,
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
	return func(db Db, ctx context.Context, t *testing.T) {
		fL := util.LiftContextualizer1V(WithTransaction[types.Unit], db, f)
		fL(ctx, t)
	}
}

// RunTestPairs executes a list of TestPair.
func RunTestPairs(db Db, ctx context.Context, t *testing.T, name string, testPairs []TestPair) {
	t.Run(name, func(t *testing.T) {
		for _, p := range testPairs {
			testFunc := func(t *testing.T) {
				p.Func(db, ctx, t)
			}
			t.Run(p.Name, testFunc)
		}
	})
}

// Parallel returns a decorated TransactionalSubtest that calls t.Parallel()
// just before executing txnlSubtest.
func Parallel(txnlSubtest TransactionalSubtest) TransactionalSubtest {
	return func(db Db, ctx context.Context, t *testing.T) {
		t.Parallel()
		txnlSubtest(db, ctx, t)
	}
}
