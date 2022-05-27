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

type Subtest func(
	ctx context.Context,
	tx pgx.Tx,
	t *testing.T,
)

type TransactionalSubtest func(
	db Db,
	ctx context.Context,
	t *testing.T,
)

type TestPair struct {
	Name string
	Func TransactionalSubtest
}

func TestWithTransaction(
	f Subtest,
) TransactionalSubtest {
	return func(db Db, ctx context.Context, t *testing.T) {
		fL := util.LiftContextualizer1V(WithTransaction[types.Unit], db, f)
		fL(ctx, t)
	}
}

func RunTestPairs(db Db, ctx context.Context, t *testing.T, testPairs []TestPair) {
	for _, p := range testPairs {
		testFunc := func(t *testing.T) {
			p.Func(db, ctx, t)
		}
		t.Run(p.Name, testFunc)
	}
}
