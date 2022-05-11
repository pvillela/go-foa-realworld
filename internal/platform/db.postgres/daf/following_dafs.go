/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
)

// FollowingCreateDaf is the instance of the DAF stereotype that
// associates a follower with a followee.
func FollowingCreateDaf(ctx context.Context, tx pgx.Tx, followerId uint, followeeId uint) error {
	sql := `
	INSERT INTO followings (follower_id, followee_id)
	VALUES ($1, $2)
	RETURNING followed_on
	`
	args := []any{followeeId, followeeId}
	_, err := tx.Exec(ctx, sql, args...)
	return errx.ErrxOf(err)
}

// FollowingDeleteDaf is the instance of the DAF stereotype that
// disaassociates a follower from a followee.
func FollowingDeleteDaf(ctx context.Context, tx pgx.Tx, followerId uint, followeeId uint) error {
	sql := `
	DELETE FROM followings
	WHERE follower_id = $1 AND followee_id = $2
	`
	_, err := tx.Exec(ctx, sql, followerId, followeeId)
	return errx.ErrxOf(err)
}
