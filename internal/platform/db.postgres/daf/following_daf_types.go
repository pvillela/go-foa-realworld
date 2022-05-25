/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// FollowingCreateDafT is the type of the stereotype instance for the DAF that
// associates a follower with a followee.
type FollowingCreateDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	followerId uint,
	followeeId uint,
) error

// FollowingGetDafT is the type of the stereotype instance for the DAF that
// retrieves an association between a follower and a followee.
type FollowingGetDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	followerId uint,
	followeeId uint,
) (model.Following, error)

// FollowingDeleteDafT is the type of the stereotype instance for the DAF that
// disaassociates a follower from a followee.
type FollowingDeleteDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	followerId uint,
	followeeId uint,
) error
