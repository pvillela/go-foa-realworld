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
	"github.com/pvillela/go-foa-realworld/experimental/arch/errx"
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
	if kind := ClassifyError(err); kind != nil {
		return kind.Make(err, "")
	}
	defer rows.Close()

	err = pgxscan.ScanOne(record, rows)
	if kind := ClassifyError(err); kind != nil {
		return kind.Make(err, "")
	}

	return nil
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
	if kind := ClassifyError(err); kind != nil {
		return nil, kind.Make(err, "")
	}
	defer rows.Close()

	var dest []R
	err = pgxscan.ScanAll(&dest, rows)
	if kind := ClassifyError(err); kind != nil {
		return nil, kind.Make(err, "")
	}

	return dest, nil
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
	DbErrRuntimeEnvironment        = errx.NewKind("DbErrRuntimeEnvironment")
	DbErrInternalAppError          = errx.NewKind("DbErrInternalAppError")
	DbErrConnectionException       = errx.NewKind("DbErrConnectionException")
	DbErrConstraintViolation       = errx.NewKind("DbErrConstraintViolation")
	DbErrUniqueViolation           = errx.NewKind("DbErrUniqueViolation", DbErrConstraintViolation)
	DbErrInsufficientResources     = errx.NewKind("DbErrInsufficientResources", DbErrRuntimeEnvironment)
	DbErrOperatorIntervention      = errx.NewKind("DbErrOperatorIntervention", DbErrRuntimeEnvironment)
	DbErrExternalSystemError       = errx.NewKind("DbErrExternalSystemError", DbErrRuntimeEnvironment)
	DbErrEngineError               = errx.NewKind("DbErrEngineError", DbErrRuntimeEnvironment)
	DbErrRecordNotFound            = errx.NewKind("DbErrRecordNotFound")
	DbErrUnexpectedMultipleRecords = errx.NewKind("DbErrUnexpectedMultipleRecords", DbErrInternalAppError)
)

// ClassifyError returns an database-related *errx.Kind that corresponds to err.
func ClassifyError(err error) *errx.Kind {
	if err == nil {
		return nil
	}

	if errX, ok := err.(errx.Errx); ok {
		return errX.Kind()
	}

	sqlState := SqlState(err)

	if sqlState == "23505" {
		return DbErrUniqueViolation
	}

	// Kludge to remove warning "second argument to errors.As should not be *error"
	var foo any = pgx.ErrNoRows
	if ok := errors.As(err, &foo); ok || pgxscan.NotFound(err) {
		return DbErrRecordNotFound
	}

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

	if strings.Contains(err.Error(), "scany") &&
		strings.Contains(err.Error(), "expected") &&
		strings.Contains(err.Error(), "row") &&
		strings.Contains(err.Error(), "got") {
		return DbErrUnexpectedMultipleRecords
	}

	return DbErrInternalAppError
}
