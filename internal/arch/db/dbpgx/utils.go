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
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"strings"
)

// SliceToWhereClause takes a slice of name-value pairs and returns a pair consisting of:
// a clause string of the form "name1 = $1 AND name2 = $2 AND ... AND nameN = $N"
// and a slice containing the values to be substituted into each $x element of the clause.
func SliceToWhereClause[V any](nvs []util.NameValuePair[string, V]) (string, []V) {
	clauseItems := make([]string, len(nvs))
	args := make([]V, len(nvs))
	for i, nv := range nvs {
		clauseItems[i] = fmt.Sprintf("%v = $%d", nv.Name, i+1)
		args[i] = nv.Value
	}
	clause := strings.Join(clauseItems, " AND ")
	return clause, args
}

// MapToWhereClause takes a map of string to an arbitrary value type and returns a pair consisting of:
// a clause string of the form "name1 = $1 AND name2 = $2 AND ... AND nameN = $N"
// and a slice containing the values to be substituted into each $x element of the clause.
func MapToWhereClause[V any](nvMap map[string]V) (string, []V) {
	nvs := make([]util.NameValuePair[string, V], len(nvMap))
	i := 0
	for n, v := range nvMap {
		nvs[i] = util.NameValuePair[string, V]{
			Name:  n,
			Value: v,
		}
	}
	return SliceToWhereClause(nvs)
}

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

func ReadMany[R any](
	ctx context.Context,
	tx pgx.Tx,
	sql string,
	args ...any,
) ([]R, error) {
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, errx.ErrxOf(err)
	}
	defer rows.Close()

	var dest []R
	err = pgxscan.ScanAll(&dest, rows)
	return dest, errx.ErrxOf(err)
}
