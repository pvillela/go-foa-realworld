/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package dbpgx

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
)

// ReadSingle reads a single record from a table.
func ReadSingle[R any, F any](
	ctx context.Context,
	tx pgx.Tx,
	tableName string,
	fieldName string,
	fieldValue F,
	record *R,
) error {
	sql := fmt.Sprintf("SELECT * FROM %v WHERE %v = $1", tableName, fieldName)
	rows, err := tx.Query(ctx, sql, fieldValue)
	if err != nil {
		return errx.ErrxOf(err)
	}
	defer rows.Close()

	err = pgxscan.ScanOne(record, rows)
	return errx.ErrxOf(err)
}

// ReadMany reads an array of records from the database. `mainSql` is the main query string,
// which is appended with optional limit and offset clauses. If limit or offset is negative
// then the corresponding clause is not appended.
func ReadMany[R any](
	ctx context.Context,
	tx pgx.Tx,
	mainSql string,
	limit int,
	offset int,
	args ...any,
) ([]R, error) {
	sql := mainSql
	if limit >= 0 {
		sql += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset >= 0 {
		sql += fmt.Sprintf(" OFFSET %d", offset)
	}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, errx.ErrxOf(err)
	}
	defer rows.Close()

	var dest []R
	err = pgxscan.ScanAll(&dest, rows)
	return dest, errx.ErrxOf(err)
}
