/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/experimental/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/experimental/arch/db/dbpgx/dbpgxtest"
	"github.com/pvillela/go-foa-realworld/experimental/arch/util"
	"github.com/pvillela/go-foa-realworld/experimental/model"
	"github.com/pvillela/go-foa-realworld/experimental/daf"
	"github.com/pvillela/go-foa-realworld/experimental/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

var followingDafsSubt = dbpgxtest.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {

	// Create followings.

	type followingSourceT struct {
		followerName string
		followeeName string
	}

	followingSources := []followingSourceT{
		{
			followerName: username1,
			followeeName: username2,
		},
		{
			followerName: username3,
			followeeName: username2,
		},
	}

	var makeFollowing = func(src followingSourceT) (uint, uint) {
		followerID := mdb.UserGetByName(src.followerName).Id
		followeeID := mdb.UserGetByName(src.followeeName).Id
		return followerID, followeeID
	}

	for _, fsrc := range followingSources {
		followerId, followeeId := makeFollowing(fsrc)
		followedOn, err := daf.FollowingCreateDaf(ctx, tx, followerId, followeeId)
		assert.NoError(t, err)
		mdb.FollowingUpsert(fsrc.followerName, fsrc.followeeName, followedOn)
	}

	// Tests

	{
		msg := "Get articles from authors followed by current user."

		currUsername := username1

		currUserId := mdb.UserGetByName(currUsername).Id

		returned, err := daf.ArticlesFeedDaf(ctx, tx, currUserId, nil, nil)
		assert.NoError(t, err)

		actual := testutil.ArticlePlusesToArticles(returned)

		expected0 := util.SliceFilter(mdb.ArticlePlusGetAll(currUsername),
			func(ap model.ArticlePlus) bool {
				return ap.Author.Following
			})
		expected := testutil.ArticlePlusesToArticles(expected0)

		assert.ElementsMatch(t, expected, actual, msg)
	}

	{
		msg := "FollowingCreateDaf - attempt to create an existing following"

		followerName := username3
		followeeName := username2

		following := mdb.FollowingGet(followerName, followeeName)
		followerId := following.FollowerId
		followeeId := following.FolloweeId

		// --> start nested transaction to avoid invalidating main transaction tx
		subTx, err := tx.Begin(ctx)
		assert.NoError(t, err)

		_, err = daf.FollowingCreateDaf(ctx, subTx, followerId, followeeId)

		retErrxKind := dbpgx.ClassifyError(err)
		expErrxKind := dbpgx.DbErrUniqueViolation

		err = subTx.Rollback(ctx)
		assert.NoError(t, err)
		// <-- rolled back nested transaction

		assert.Equal(t, expErrxKind, retErrxKind, msg)
	}

	{
		msg := "FollowingDeleteDaf -- delete an existing following."

		followerName := username3
		followeeName := username2

		following := mdb.FollowingGet(followerName, followeeName)
		followerId := following.FollowerId
		followeeId := following.FolloweeId

		err := daf.FollowingDeleteDaf(ctx, tx, followerId, followeeId)
		assert.NoError(t, err)

		mdb.FollowingDelete(followerName, followeeName)

		returned, err := daf.ArticlesFeedDaf(ctx, tx, followerId, nil, nil)
		assert.NoError(t, err)

		actual := testutil.ArticlePlusesToArticles(returned)

		expected0 := util.SliceFilter(mdb.ArticlePlusGetAll(followerName),
			func(ap model.ArticlePlus) bool {
				return ap.Author.Following
			})
		expected := testutil.ArticlePlusesToArticles(expected0)

		assert.ElementsMatch(t, expected, actual, msg)
	}

	{
		msg := "FollowingDeleteDaf - attempt to delete a nonexistent following"

		followerName := username3
		followeeName := username2

		following := mdb.FollowingGet(followerName, followeeName)
		followerId := following.FollowerId
		followeeId := following.FolloweeId

		err := daf.FollowingDeleteDaf(ctx, tx, followerId, followeeId)
		retErrxKind := dbpgx.ClassifyError(err)
		expErrxKind := dbpgx.DbErrRecordNotFound

		assert.Equal(t, expErrxKind, retErrxKind, msg)
	}
})
