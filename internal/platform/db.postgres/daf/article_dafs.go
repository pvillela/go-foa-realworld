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
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"strings"
)

const (
	limitDefault  = 20
	offsetDefault = 0
)

// ArticleCreateDaf implements a stereotype instance of type
// ArticleCreateDafT.
var ArticleCreateDaf ArticleCreateDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	article *model.Article,
) error {
	sql := `
	INSERT INTO articles (author_id, title, slug, description, body, favorites_count) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at
	`
	args := []any{
		article.Author.Id,
		article.Title,
		article.Slug,
		article.Description,
		article.Body,
		article.FavoritesCount,
	}

	row := tx.QueryRow(ctx, sql, args...)
	err := row.Scan(&article.Id, &article.CreatedAt, &article.UpdatedAt)
	return errx.ErrxOf(err)
}

// ArticleGetBySlugDaf implements a stereotype instance of type
// ArticleGetBySlugDafT.
var ArticleGetBySlugDaf ArticleGetBySlugDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	currUserId uint,
	slug string,
) (model.ArticlePlus, error) {
	// This implementation uses the readArticles helper function. Although a simpler direct
	// implementation is clearly possible, the required SQL query and data mapping logic are
	// sufficiently complex that it is best to deal with those things only once.

	// See selectJoin query string in function readArticles
	whereTuples := []util.Tuple2[string, any]{
		{"a.slug = $%d", slug},
		{"ufa2.username = $%d", ""}, // use invalid username to force null join
	}

	where, whereArgs := whereClauseFromTuples(whereTuples)

	articlePluses, _, err := readArticles(ctx, tx, currUserId, where, nil, nil, whereArgs...)
	if err != nil {
		return model.ArticlePlus{}, err
	}

	if len(articlePluses) == 0 {
		return model.ArticlePlus{}, fs.ErrArticleSlugNotFound.Make(nil, slug)
	}
	if len(articlePluses) > 1 {
		util.PanicOnError(errx.NewErrx(nil,
			fmt.Sprintf("Found multiple articles with same slug '%v'", slug)))
	}

	return articlePluses[0], nil
}

// ArticleUpdateDaf implements a stereotype instance of type
// ArticleUpdateDafT.
var ArticleUpdateDaf ArticleUpdateDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	article *model.Article,
) error {
	sql := `
	UPDATE articles 
	SET title = $1, description = $2, body = $3, updated_at = NOW() 
	WHERE id = $4 AND updated_at = $5
	RETURNING updated_at
	`
	args := []interface{}{
		article.Title,
		article.Description,
		article.Body,
		article.Id,
		article.UpdatedAt,
	}
	row := tx.QueryRow(ctx, sql, args...)
	err := row.Scan(&article.UpdatedAt)
	return errx.ErrxOf(err)
}

// ArticleDeleteDaf implements a stereotype instance of type
// ArticleDeleteDafT.
var ArticleDeleteDaf ArticleDeleteDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	slug string,
) error {
	sql := `
	DELETE FROM articles
	WHERE slug = $1
	`
	_, err := tx.Exec(ctx, sql, slug)
	return errx.ErrxOf(err)
}

// ArticlesFeedDaf implements a stereotype instance of type
// ArticlesFeedDafT.
var ArticlesFeedDaf ArticlesFeedDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	currUserId uint,
	optLimit *int,
	optOffset *int,
) ([]model.ArticlePlus, error) {
	// See selectJoin query string in function readArticles
	whereTuples := []util.Tuple2[string, any]{
		{"fo.follower_id = $%d", currUserId},
		{"ufa2.username = $%d", ""}, // use invalid username to force null join
	}

	where, whereArgs := whereClauseFromTuples(whereTuples)

	articlePluses, _, err := readArticles(ctx, tx, currUserId, where, optLimit, optOffset, whereArgs...)
	if err != nil {
		return nil, err
	}

	return articlePluses, nil
}

// ArticlesListDafT implements a stereotype instance of type
// ArticlesListDafTT.
var ArticlesListDaf ArticlesListDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	currUserId uint,
	criteria model.ArticleCriteria,
) ([]model.ArticlePlus, error) {
	// See selectJoin query string in function readArticles
	whereTuples := make([]util.Tuple2[string, any], 0)
	if v := criteria.Tag; v != nil {
		whereTuples = append(whereTuples, util.Tuple2[string, any]{"t.name = $%d", *v})
	}
	if v := criteria.Author; v != nil {
		whereTuples = append(whereTuples, util.Tuple2[string, any]{"ua.username = $%d", *v})
	}
	if v := criteria.FavoritedBy; v != nil {
		whereTuples = append(whereTuples, util.Tuple2[string, any]{"ufa2.username = $%d", *v})
	} else {
		// use invalid username to force null join
		whereTuples = append(whereTuples, util.Tuple2[string, any]{"ufa2.username = $%d", ""})
	}

	where, whereArgs := whereClauseFromTuples(whereTuples)

	articlePluses, _, err := readArticles(ctx, tx, currUserId, where, criteria.Limit, criteria.Offset,
		whereArgs...)
	if err != nil {
		return nil, err
	}

	return articlePluses, nil
}

/////////////////////
// Helper functions

func whereClauseFromTuples(whereTuples []util.Tuple2[string, any]) (where string, whereArgs []any) {
	wheres := make([]string, len(whereTuples))
	whereArgs = make([]any, len(whereTuples))
	for i, nv := range whereTuples {
		wheres[i] = fmt.Sprintf(nv.X1, i+2)
		whereArgs[i] = nv.X2
	}
	where = strings.Join(wheres, " AND ")
	if len(wheres) != 0 {
		where = " WHERE " + where
	}
	return where, whereArgs
}

// readArticles queries for articles, returning slice of model.ArticlePlus and an error.
// The first slice contains model.ArticlePlus items.
// The second slice contains, for each item, a slice of the user IDs of all users that have
// favorited the article with the same index in the first slice. This second slice is not
// very useful, except maybe for debugging purposes.
//
// currUserId is the user ID of the currently logged-in user.
//
// where is a WHERE clause string that filters the query. The argument placeholders in `where`
// must start with $2 as $1 is reserved for `currUserId`. The query, implemented inside this
// function, provides a denormalized holistic view of articles by joining the tables: articles,
// users (twice), favorites (twice), followings, and tags. The double joins are required to
// support the where clauses from the DAFs that call this function.
//
// optLimit and optOffset are optional limit and offset parameters for the result set.
//
// whereArgs are arguments for the where clause.
func readArticles(
	ctx context.Context,
	tx pgx.Tx,
	currUserId uint,
	where string,
	optLimit *int,
	optOffset *int,
	whereArgs ...any,
) ([]model.ArticlePlus, [][]uint, error) {
	// Construct SQL
	selectJoin := `
	SELECT a.*, ua.username, ua.bio, ua.image, fa1.user_id as favorited, fa2.user_id as favorites_user_id, 
		fo.follower_id as following, t.name 
		FROM articles a
	LEFT JOIN users ua ON a.author_id = ua.id
	LEFT JOIN users ufa2 ON fa2.user_id = ufa2.id -- product effect, same as fa2
	LEFT JOIN favorites fa1 ON a.id = fa1.article_id AND $1 = fa1.user_id -- no product effect
	LEFT JOIN favorites fa2 ON a.id = fa2.article_id -- product effect
	LEFT JOIN followings fo ON fo.follower_id = $1 AND a.author_id = fo.followee_id -- no product effect
	LEFT JOIN article_tags at ON a.id = at.article_id -- product effect
	LEFT JOIN tags t ON at.tag_id = t.id
	`
	orderBy := `
	ORDER BY a.created_at DESC, t.name, favorites_user_id
	`
	sql := selectJoin + where + orderBy
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

	// Retrieve rows
	args := append([]any{currUserId}, whereArgs...)
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, nil, errx.ErrxOf(err)
	}
	defer rows.Close()

	// Data structures to receive data
	type resultT struct {
		articlePlus      model.ArticlePlus
		favoritesUserIds []uint
	}
	results := []resultT{}
	tags := []string{}
	favoritesUserIds := []uint{}
	var currRecord struct {
		article         model.Article `db:""`
		profile         model.Profile `db:""`
		following       *uint
		favorited       *uint
		tag             *string `db:"name"`
		favoritesUserId *uint
	}
	var prevRecord struct {
		article         model.Article `db:""`
		profile         model.Profile `db:""`
		following       *uint
		favorited       *uint
		tag             *string `db:"name"`
		favoritesUserId *uint
	}

	// Functions to put data in data structures

	fresh := true // true when first encountering a new article

	articleChanged := func() bool {
		return !fresh && currRecord.article.Id != prevRecord.article.Id
	}

	tagChanged := func() bool {
		return fresh ||
			(currRecord.tag == nil && prevRecord.tag != nil) ||
			(currRecord.tag != nil && prevRecord.tag == nil) ||
			(*currRecord.tag != *prevRecord.tag)
	}

	favoritesUserIdChanged := func() bool {
		return fresh ||
			(currRecord.favoritesUserId == nil && prevRecord.favoritesUserId != nil) ||
			(currRecord.favoritesUserId != nil && prevRecord.favoritesUserId == nil) ||
			(*currRecord.favoritesUserId != *prevRecord.favoritesUserId)
	}

	// Move data into results slice
	updateResults := func() {
		ra := &prevRecord.article
		rp := &prevRecord.profile
		if prevRecord.following != nil {
			rp.Following = true
		}
		favorited := false
		if prevRecord.favorited != nil {
			favorited = true
		}
		articlePlus := model.ArticlePlus{
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
		result := resultT{
			articlePlus:      articlePlus,
			favoritesUserIds: favoritesUserIds,
		}
		results = append(results, result)
		tags = []string{}
		favoritesUserIds = []uint{}
	}

	// Main processing loop
	for rows.Next() {
		err = pgxscan.ScanRow(&currRecord, rows)
		if err != nil {
			return nil, nil, errx.ErrxOf(err)
		}

		if articleChanged() {
			updateResults()
			fresh = true
		}

		if tagChanged() && currRecord.tag != nil {
			tags = append(tags, *currRecord.tag)
		}

		if favoritesUserIdChanged() && currRecord.favoritesUserId != nil {
			favoritesUserIds = append(favoritesUserIds, *currRecord.favoritesUserId)
		}

		prevRecord = currRecord
		fresh = false
	}
	if !fresh {
		updateResults()
	}

	articlePluses := make([]model.ArticlePlus, len(results))
	favoritesUserIdSliceSlice := make([][]uint, len(results))
	for i, r := range results {
		articlePluses[i] = r.articlePlus
		favoritesUserIdSliceSlice[i] = r.favoritesUserIds
	}

	return articlePluses, favoritesUserIdSliceSlice, nil
}
