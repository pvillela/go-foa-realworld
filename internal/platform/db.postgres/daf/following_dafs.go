/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/bf"
)

// FollowingCreateDafI is the instance of the DAF stereotype that
// associates a follower with a followee.
var FollowingCreateDafI FollowingCreateDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	followerId uint,
	followeeId uint,
) error {
	sql := `
	INSERT INTO followings (follower_id, followee_id)
	VALUES ($1, $2)
	RETURNING followed_on
	`
	args := []any{followerId, followeeId}

	_, err := tx.Exec(ctx, sql, args...)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind == dbpgx.DbErrUniqueViolation {
			return kind.Make(err, bf.ErrMsgUserAlreadyFollowed, followeeId)
		}
		return kind.Make(err, "")
	}

	return nil
}

// FollowingDeleteDafI is the instance of the DAF stereotype that
// disaassociates a follower from a followee.
var FollowingDeleteDafI FollowingDeleteDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	followerId uint,
	followeeId uint,
) error {
	sql := `
	DELETE FROM followings
	WHERE follower_id = $1 AND followee_id = $2
	`
	args := []any{followerId, followeeId}
	c, err := tx.Exec(ctx, sql, args...)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		return kind.Make(err, "")
	}
	if c.RowsAffected() == 0 {
		return dbpgx.DbErrRecordNotFound.Make(nil, bf.ErrMsgUserWasNotFollowed, followeeId)
	}

	return nil
}
