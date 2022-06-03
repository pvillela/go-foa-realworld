/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"time"
)

// FollowingCreateDafI is the instance of the DAF stereotype that
// associates a follower with a followee.
var FollowingCreateDafI FollowingCreateDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	followerId uint,
	followeeId uint,
) (time.Time, error) {
	sql := `
	INSERT INTO followings (follower_id, followee_id)
	VALUES ($1, $2)
	RETURNING followed_on
	`
	args := []any{followerId, followeeId}

	var followedOn time.Time
	row := tx.QueryRow(ctx, sql, args...)
	err := row.Scan(&followedOn)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind == dbpgx.DbErrRecordNotFound {
			return time.Time{},
				dbpgx.DbErrUniqueViolation.Make(err, bf.ErrMsgUserAlreadyFollowed, followeeId)
		}
		return time.Time{}, kind.Make(err, "")
	}

	return followedOn, nil
}

// FollowingGetDafI implements a stereotype instance of type
// FollowingGetDafT.
// Returns the association record if found or a zero model.Following otherwise.
var FollowingGetDafI FollowingGetDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	followerId uint,
	followeeId uint,
) (model.Following, error) {
	var zero model.Following

	sql := `
	SELECT * FROM followings 
	WHERE follower_id = $1 AND followee_id = $2
	`
	args := []any{
		followerId,
		followeeId,
	}

	rows, err := tx.Query(ctx, sql, args...)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		return zero, kind.Make(err, "")
	}
	defer rows.Close()

	var following model.Following
	err = pgxscan.ScanOne(&following, rows)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind != dbpgx.DbErrRecordNotFound {
			return zero, kind.Make(err, "")
		}
		// Otherwise, the association was not found and the zero value will be returned.
	}

	return following, nil
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
