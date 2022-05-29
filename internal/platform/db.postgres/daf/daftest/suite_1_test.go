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

func TestSuite1(t *testing.T) {
	mainTxnSubtest := func(db dbpgx.Db, ctx context.Context, t *testing.T) {
		dbpgx.RunTestPairs(db, ctx, t, "user_and_article_parallel", []dbpgx.TestPair{
			{Name: "userDafsSubt", Func: dbpgx.Parallel(UserDafsSubt)},
			{Name: "articleDafsSubt", Func: dbpgx.Parallel(articleDafsSubt)},
		})
		dbpgx.RunTestPairs(db, ctx, t, "other_dafs", []dbpgx.TestPair{
			{Name: "commentDafsSubt", Func: commentDafsSubt},
		})
	}

	dafTester(t, "DafTests", mainTxnSubtest)
}
