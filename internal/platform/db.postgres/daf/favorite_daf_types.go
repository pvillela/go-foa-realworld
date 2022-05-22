/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4"
)

// FavoriteCreateDafT is the instance of the DAF stereotype that
// associates an article with a user that favors it.
type FavoriteCreateDafT = func(ctx context.Context, tx pgx.Tx, articleId uint, userId uint) error

// FavoriteDeleteDafT is the instance of the DAF stereotype that
// disaassociates an article from a user that favors it.
// Returns the number of rows affected, which can be 0 or 1.
type FavoriteDeleteDafT = func(ctx context.Context, tx pgx.Tx, articleId uint, userId uint) (int, error)
