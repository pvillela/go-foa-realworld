/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfltest

import (
	"context"
	"testing"

	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx/dbpgxtest"
	"github.com/pvillela/go-foa-realworld/internal/testutil"
	"github.com/sirupsen/logrus"
)

var connStr = "postgres://testuser:testpassword@localhost:9999/testdb?sslmode=disable"

func TestSflSuite(t *testing.T) {

	logrus.SetLevel(logrus.DebugLevel)

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableQuote: true,
	})

	txnlSubtest := func(db dbpgx.Db, ctx context.Context, t *testing.T) {
		testPairs := []dbpgxtest.TestPair{
			{Name: "userRegisterSflSubt", Func: userRegisterSflSubt},
			{Name: "userAuthenticateSflSubt", Func: userAuthenticateSflSubt},
			{Name: "userFollowSflSubt", Func: userFollowSflSubt},
			{Name: "userGetCurrentSflSubt", Func: userGetCurrentSflSubt},
			{Name: "userUnfollowSflSubt", Func: userUnfollowSflSubt},
			// Undo above unfollow action to enable subsequent tests
			{Name: "userFollowSflSubt", Func: userFollowSflSubt},
			{Name: "userUpdateSflSubt", Func: userUpdateSflSubt},
			{Name: "articleCreateSflSubt", Func: articleCreateSflSubt},
			{Name: "articleDeleteSflSubt", Func: articleDeleteSflSubt},
			{Name: "articleFavoriteSflSubt", Func: articleFavoriteSflSubt},
			{Name: "articleGetSflSubt", Func: articleGetSflSubt},
			{Name: "articleUnfavoriteSflSubt", Func: articleUnfavoriteSflSubt},
			{Name: "articleUpdateSflSubt", Func: articleUpdateSflSubt},
			{Name: "articlesFeedSflSubt", Func: articlesFeedSflSubt},
			{Name: "articlesListSflSubt", Func: articlesListSflSubt},
			{Name: "commentAddSflSubt", Func: commentAddSflSubt},
			//{Name: "commentSflsSubt", Func: commentSflsSubt},
			//{Name: "profileSflsSubt", Func: profileSflsSubt},
			//{Name: "tagSflsSubt", Func: tagSflsSubt},
		}

		dbpgxtest.RunTestPairs(db, ctx, t, "sequential", testPairs)
	}

	dbpgxtest.DbTester(t, txnlSubtest, connStr, testutil.CleanupAllTables)
}
