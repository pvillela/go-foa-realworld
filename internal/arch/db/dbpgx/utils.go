/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package dbpgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"strings"
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

// SqlState returns the the pgx SQLState() of err if err is wraps a *pgconn.PgError,
// returns the empty string otherwise.
func SqlState(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.SQLState()
	}
	return ""
}

var (
	DbErrRuntimeEnvironment        = errx.NewKind("database error runtime environment")
	DbErrInternalAppError          = errx.NewKind("database internal application error")
	DbErrConnectionException       = errx.NewKind("database connection exception")
	DbErrConstraintViolation       = errx.NewKind("database error constraint violation")
	DbErrUniqueViolation           = errx.NewKind("database error unique violation", DbErrConstraintViolation)
	DbErrInsufficientResources     = errx.NewKind("database error insufficient resources", DbErrRuntimeEnvironment)
	DbErrOperatorIntervention      = errx.NewKind("database error operator intervention", DbErrRuntimeEnvironment)
	DbErrExternalSystemError       = errx.NewKind("database external system error", DbErrRuntimeEnvironment)
	DbErrEngineError               = errx.NewKind("database engine error", DbErrRuntimeEnvironment)
	DbErrRecordNotFound            = errx.NewKind("database error record not found")
	DbErrUnexpectedMultipleRecords = errx.NewKind("database error unexpected multiple records")
)

// ClassifyError returns an database-related *errx.Kind that corresponds to err.
func ClassifyError(err error) *errx.Kind {
	if ok := errors.As(err, &pgx.ErrNoRows); ok {
		return DbErrRecordNotFound
	}

	sqlState := SqlState(err)
	prefix := sqlState[:2]
	switch prefix {
	case "08":
		return DbErrConnectionException
	case "23":
		return DbErrConstraintViolation
	case "53":
		return DbErrInsufficientResources
	case "57":
		return DbErrOperatorIntervention
	case "58":
		return DbErrExternalSystemError
	case "XX":
		return DbErrEngineError
	}

	if sqlState == "23505" {
		return DbErrUniqueViolation
	}

	if strings.Contains(err.Error(), "scany") &&
		strings.Contains(err.Error(), "expected") &&
		strings.Contains(err.Error(), "row") &&
		strings.Contains(err.Error(), "got") {
		return DbErrUnexpectedMultipleRecords
	}

	return DbErrInternalAppError
}
