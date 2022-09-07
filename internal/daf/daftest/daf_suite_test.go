/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx/dbpgxtest"
	"github.com/pvillela/go-foa-realworld/internal/testutil"
	"github.com/sirupsen/logrus"
	"testing"
)

var connStr = "postgres://testuser:testpassword@localhost:9999/testdb?sslmode=disable"

func TestDafSuite(t *testing.T) {

	logrus.SetLevel(logrus.DebugLevel)

	txnlSubtest := func(db dbpgx.Db, ctx context.Context, t *testing.T) {
		testPairs := []dbpgxtest.TestPair{
			{Name: "userDafsSubt", Func: userDafsSubt},
			{Name: "articleDafsSubt", Func: articleDafsSubt},
			{Name: "commentDafsSubt", Func: commentDafsSubt},
			{Name: "favoriteDafsSubt", Func: favoriteDafsSubt},
			{Name: "followingDafsSubt", Func: followingDafsSubt},
			{Name: "tagDafsSubt", Func: tagDafsSubt},
		}

		dbpgxtest.RunTestPairs(db, ctx, t, "sequential", testPairs)
	}

	dbpgxtest.DbTester(t, txnlSubtest, connStr, testutil.CleanupAllTables)
}
