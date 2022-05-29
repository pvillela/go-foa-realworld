/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"testing"
)

func TestDafSuite(t *testing.T) {
	txnlSubtest := func(db dbpgx.Db, ctx context.Context, t *testing.T) {
		testPairs := []dbpgx.TestPair{
			{Name: "userDafsSubt", Func: UserDafsSubt},
			{Name: "articleDafsSubt", Func: articleDafsSubt},
			{Name: "commentDafsSubt", Func: commentDafsSubt},
		}

		// It is OK to run the tests in parallel because each executes in a serializable transaction.
		parallelTestPairs := util.SliceMap(testPairs, func(tp dbpgx.TestPair) dbpgx.TestPair {
			return dbpgx.TestPair{
				Name: tp.Name,
				Func: dbpgx.Parallel(tp.Func),
			}
		})
		util.Ignore(parallelTestPairs)

		dbpgx.RunTestPairs(db, ctx, t, "sequential", testPairs)
		dbpgx.RunTestPairs(db, ctx, t, "parallel", parallelTestPairs)
	}

	dafTester(t, txnlSubtest)
}
