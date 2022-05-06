/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	log "github.com/sirupsen/logrus"
	"strings"
)

/////////////////////
// Helper functions for Article DAFs

const (
	limitDefault  = 20
	offsetDefault = 0
)

// whereClauseFromTuples returns a parameterized WHERE clause and corresponding arguments slice from
// a slice of tuples.
// Each tuple in the input slice contributes to the where clause.
// The output arguments slice contains the second components of the tuples, skiping any nil values.
// The placeholders in the WHERE clause are numbered starting with initialIndex.
func whereClauseFromTuples(
	initialIndex int,
	whereTuples []util.Tuple2[string, any],
) (where string, whereArgs []any) {
	wheres := make([]string, len(whereTuples))
	nullableWhereArgs := make([]any, len(whereTuples))
	for i, nv := range whereTuples {
		wheres[i] = nv.X1
		nullableWhereArgs[i] = nv.X2
	}
	return whereClauseFromSlices(initialIndex, wheres, nullableWhereArgs)
}

// mainArticlePlusQuery defines a denormalized "view" that retrieves all model.ArticlePlus objects.
// This "view" can be constrained by appending additional JOIN and/or WHERE clauses to it.
var mainArticlePlusQuery = `
	SELECT a.*, ua.username, ua.bio, ua.image, fa.user_id as favorited, fo.follower_id as following, t.name 
		FROM articles a
	LEFT JOIN users ua ON a.author_id = ua.id
	LEFT JOIN favorites fa ON fa.article_id = a.id AND fa.user_id = $1 -- at most one
	LEFT JOIN followings fo ON fo.follower_id = $1 AND fo.followee_id = a.author_id -- at most one
	LEFT JOIN article_tags at ON at.article_id = a.id -- product effect
	LEFT JOIN tags t ON t.id = at.tag_id 
`

// readArticles retrieves articles, returning a slice of model.ArticlePlus and an error.
// It is used by the DAFs that retrieve articles.
//
// Parameters other than the obvious ctx and tx:
//
// currUserId is the user ID of the currently logged-in user.
//
// additionaSql contains JOINs and/OR a WHERE clause that are appended to mainArticlePlusQuery to
// filter the query.
// The argument placeholders in additionalSql must start with $2 as $1 is reserved for currUserId.
//
// optLimit and optOffset are optional limit and offset parameters for the result set.
//
// additionalArgs are arguments for additionalSql.
func readArticles(
	ctx context.Context,
	tx pgx.Tx,
	currUserId uint,
	additionalSql string,
	optLimit *int,
	optOffset *int,
	additionalArgs ...any,
) ([]model.ArticlePlus, error) {
	// Construct SQL
	orderBy := `
	ORDER BY a.created_at DESC, t.name
	`
	sql := mainArticlePlusQuery + additionalSql + orderBy
	limit := limitDefault
	if optLimit != nil {
		limit = *optLimit
		sql += fmt.Sprintf(" LIMIT %d", limit)
	}
	offset := offsetDefault
	if optOffset != nil {
		offset = *optOffset
		sql += fmt.Sprintf(" OFFSET %d", offset)
	}

	log.Debug("Full sql: ", sql)

	// Retrieve rows
	args := append([]any{currUserId}, additionalArgs...)
	log.Debug("args: ", args)
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, errx.ErrxOf(err)
	}
	defer rows.Close()

	// Data structures to receive data
	type recordT struct {
		Id        uint
		Article   model.Article `db:""`
		Profile   model.Profile `db:""`
		Following *uint
		Favorited *uint
		Tag       *string `db:"name"`
	}
	var results []model.ArticlePlus
	var tags []string
	var currRecord recordT
	var prevRecord recordT

	// Functions to put data in data structures

	fresh := true // true when first encountering a new article

	articleChanged := func() bool {
		return !fresh && currRecord.Id != prevRecord.Id
	}

	tagChanged := func() bool {
		return fresh ||
			(currRecord.Tag == nil && prevRecord.Tag != nil) ||
			(currRecord.Tag != nil && prevRecord.Tag == nil) ||
			(currRecord.Tag != nil && prevRecord.Tag != nil &&
				*currRecord.Tag != *prevRecord.Tag)
	}

	// Move data into results slice
	updateResults := func() {
		ra := &prevRecord.Article
		rp := &prevRecord.Profile
		if prevRecord.Following != nil {
			rp.Following = true
		}
		favorited := false
		if prevRecord.Favorited != nil {
			favorited = true
		}
		result := model.ArticlePlus{
			Slug: ra.Slug,
			Author: model.Profile{
				Username:  rp.Username,
				Bio:       rp.Bio,
				Image:     rp.Image,
				Following: rp.Following,
			},
			Title:          ra.Title,
			Description:    ra.Description,
			Body:           ra.Body,
			TagList:        tags,
			CreatedAt:      ra.CreatedAt,
			UpdatedAt:      ra.UpdatedAt,
			Favorited:      favorited,
			FavoritesCount: ra.FavoritesCount,
		}
		results = append(results, result)
		tags = []string{}
	}

	// Main processing loop
	for rows.Next() {
		err = pgxscan.ScanRow(&currRecord, rows)
		if err != nil {
			return nil, errx.ErrxOf(err)
		}

		if articleChanged() {
			updateResults()
			fresh = true
		}

		if tagChanged() && currRecord.Tag != nil {
			tags = append(tags, *currRecord.Tag)
		}

		prevRecord = currRecord
		fresh = false
	}
	if !fresh {
		updateResults()
	}

	return results, nil
}

// Helper to whereClauseFromTuples; do not use it directly.
func whereClauseFromSlices(
	initialIndex int,
	wheres []string,
	nullableWhereArgs []any,
) (where string, whereArgs []any) {
	if len(wheres) == 0 {
		return "", nullableWhereArgs
	}

	wheresWithParams := make([]string, len(wheres))
	whereArgs = make([]any, 0, len(wheres))
	idx := initialIndex
	for i, w := range wheres {
		if nullableWhereArgs[i] != nil {
			wheresWithParams[i] = fmt.Sprintf(w, idx)
			whereArgs = append(whereArgs, nullableWhereArgs[i])
			idx++
		} else {
			wheresWithParams[i] = wheres[i]
		}
	}

	where = "WHERE " + strings.Join(wheresWithParams, " AND ")
	return where, whereArgs
}
