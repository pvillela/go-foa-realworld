/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"testing"
)

func TestDafSuite(t *testing.T) {
	txnlSubtest := func(db dbpgx.Db, ctx context.Context, t *testing.T) {
		testPairs := []dbpgx.TestPair{
			//{Name: "userDafsSubt", Func: userDafsSubt},
			//{Name: "articleDafsSubt", Func: articleDafsSubt},
			//{Name: "commentDafsSubt", Func: dbpgx.Parallel(commentDafsSubt)},
			//{Name: "favoriteDafsSubt", Func: dbpgx.Parallel(favoriteDafsSubt)},
		}

		dbpgx.RunTestPairs(db, ctx, t, "sequential", testPairs)
	}

	dafTester(t, txnlSubtest)
}
