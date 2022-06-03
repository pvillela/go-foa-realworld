/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/stretchr/testify/assert"
	"testing"
)

var followingDafsSubt = dbpgx.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {

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
		followedOn, err := daf.FollowingCreateDafI(ctx, tx, followerId, followeeId)
		errx.PanicOnError(err)
		mdb.FollowingUpsert(fsrc.followerName, fsrc.followeeName, followedOn)
	}

	// Tests

	{
		msg := "Get articles from authors followed by current user."

		currUsername := username1

		currUserId := mdb.UserGetByName(currUsername).Id

		returned, err := daf.ArticlesFeedDafI(ctx, tx, currUserId, nil, nil)
		errx.PanicOnError(err)

		actual := ArticlePlusesToArticles(returned)

		expected0 := util.SliceFilter(mdb.ArticlePlusGetAll(currUsername),
			func(ap model.ArticlePlus) bool {
				return ap.Author.Following
			})
		expected := ArticlePlusesToArticles(expected0)

		assert.ElementsMatch(t, expected, actual, msg)
	}

	//{
	//	currUserId := users[2].Id
	//
	//	articlePluses, err := daf.ArticlesFeedDafI(ctx, tx, currUserId, nil, nil)
	//	errx.PanicOnError(err)
	//	fmt.Println("\nArticlesFeedDaf:", articlePluses, "\n")
	//}
	//
	//{
	//	err := daf.FollowingDeleteDafI(ctx, tx, followings[1].FollowerID, followings[1].FolloweeID)
	//	errx.PanicOnError(err)
	//
	//	currUserId := users[2].Id
	//
	//	articlePluses, err := daf.ArticlesFeedDafI(ctx, tx, currUserId, nil, nil)
	//	errx.PanicOnError(err)
	//	fmt.Println("\nArticlesFeedDaf:", articlePluses, "\n")
	//}
})
