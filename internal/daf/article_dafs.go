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
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/rpc"
	"time"
)

/////////////////////
// Types

// ArticleCreateDafT is the type of the stereotype instance for the DAF that
// creates an article.
type ArticleCreateDafT = func(ctx context.Context, tx pgx.Tx, article *model.Article) error

// ArticleGetBySlugDafT is the type of the stereotype instance for the DAF that
// retrieves an article by slug.
type ArticleGetBySlugDafT = func(ctx context.Context, tx pgx.Tx, currUserId uint, slug string) (model.ArticlePlus, error)

// ArticleUpdateDafT is the type of the stereotype instance for the DAF that
// updates an article.
type ArticleUpdateDafT = func(ctx context.Context, tx pgx.Tx, article *model.Article) error

// ArticleAdjustFavoritesCountDafT is the type of the stereotype instance for the DAF that
// increments the favorites count of an article.
type ArticleAdjustFavoritesCountDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	slug string,
	delta int,
) (int, time.Time, error)

// ArticleDeleteDafT is the type of the stereotype instance for the DAF that
// deletes an article.
type ArticleDeleteDafT = func(ctx context.Context, tx pgx.Tx, slug string) error

// ArticlesFeedDafT is the type of the stereotype instance for the DAF that
// queries for all articles authored by other users followed by
// the current user, with optional limit and offset pagination parameters.
type ArticlesFeedDafT = func(ctx context.Context, tx pgx.Tx, currUserId uint, optLimit *int, optOffset *int) ([]model.ArticlePlus, error)

// ArticlesListDafT is the type of the stereotype instance for the DAF that
// retrieve recent articles based on a set of query parameters.
type ArticlesListDafT = func(ctx context.Context, tx pgx.Tx, currUserId uint, criteria rpc.ArticleCriteria) ([]model.ArticlePlus, error)

/////////////////////
// DAFS

// ArticleCreateDaf implements a stereotype instance of type
// ArticleCreateDafT.
var ArticleCreateDaf ArticleCreateDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	article *model.Article,
) error {
	sql := `
	INSERT INTO articles (author_id, title, slug, description, body, favorites_count) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id, created_at, updated_at
	`
	args := []any{
		article.AuthorId,
		article.Title,
		article.Slug,
		article.Description,
		article.Body,
		article.FavoritesCount,
	}

	row := tx.QueryRow(ctx, sql, args...)
	err := row.Scan(&article.Id, &article.CreatedAt, &article.UpdatedAt)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		if kind == dbpgx.DbErrUniqueViolation {
			return kind.Make(err, bf.ErrMsgDuplicateArticleSlug, article.Slug)
		}
		return kind.Make(err, "")
	}
	return nil
}

// ArticleGetBySlugDaf implements a stereotype instance of type
// ArticleGetBySlugDafT.
var ArticleGetBySlugDaf ArticleGetBySlugDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	currUserId uint,
	slug string,
) (model.ArticlePlus, error) {
	// See mainArticlePlusQuery
	whereTuples := []types.Tuple2[string, any]{
		{"a.slug = $%d", slug},
	}
	where, whereArgs := whereClauseFromTuples(2, whereTuples)

	results, err := readArticles(ctx, tx, currUserId, where, nil, nil, whereArgs...)
	if err != nil {
		return model.ArticlePlus{}, err
	}
	if len(results) == 0 {
		return model.ArticlePlus{},
			dbpgx.DbErrRecordNotFound.Make(nil, bf.ErrMsgArticleSlugNotFound, slug)
	}
	if len(results) > 1 {
		return model.ArticlePlus{},
			dbpgx.DbErrUnexpectedMultipleRecords.Make(nil,
				"Found multiple articles with same slug '%v'", slug)
	}

	return results[0], nil
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
	SET title = $1, description = $2, body = $3, updated_at = clock_timestamp() 
	WHERE slug = $4 AND updated_at = $5
	RETURNING updated_at
	`
	args := []any{
		article.Title,
		article.Description,
		article.Body,
		article.Slug,
		article.UpdatedAt,
	}
	row := tx.QueryRow(ctx, sql, args...)
	if err := row.Scan(&article.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			err = dbpgx.DbErrRecordNotFound.Make(nil, bf.ErrMsgArticleSlugNotFound, article.Slug)
		}
		return dbpgx.ClassifyError(err).Make(err, "")
	}

	return nil
}

// ArticleAdjustFavoritesCountDaf implements a stereotype instance of type
// ArticleAdjustFavoritesCountDafT.
var ArticleAdjustFavoritesCountDaf ArticleAdjustFavoritesCountDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	slug string,
	delta int,
) (int, time.Time, error) {
	sql := `
	UPDATE articles 
	SET favorites_count = favorites_count + $2, updated_at = clock_timestamp() 
	WHERE slug = $1
	RETURNING favorites_count, updated_at
	`
	args := []any{
		slug,
		delta,
	}

	var favoritesCount int
	var updatedAt time.Time

	row := tx.QueryRow(ctx, sql, args...)
	if err := row.Scan(&favoritesCount, &updatedAt); err != nil {
		if err == pgx.ErrNoRows {
			err = dbpgx.DbErrRecordNotFound.Make(nil, bf.ErrMsgArticleSlugNotFound, slug)
		}
		return 0, time.Time{}, dbpgx.ClassifyError(err).Make(err, "")
	}

	return favoritesCount, updatedAt, nil
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
	c, err := tx.Exec(ctx, sql, slug)
	if kind := dbpgx.ClassifyError(err); kind != nil {
		return kind.Make(err, "")
	}
	if c.RowsAffected() == 0 {
		return dbpgx.DbErrRecordNotFound.Make(nil, bf.ErrMsgArticleSlugNotFound, slug)
	}

	return nil
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
	// See mainArticlePlusQuery
	whereTuples := []types.Tuple2[string, any]{
		{"fo.follower_id = $%d", currUserId},
	}
	where, whereArgs := whereClauseFromTuples(2, whereTuples)

	results, err := readArticles(ctx, tx, currUserId, where, optLimit, optOffset, whereArgs...)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// ArticlesListDaf implements a stereotype instance of type
// ArticlesListDafT.
var ArticlesListDaf ArticlesListDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	currUserId uint,
	criteria rpc.ArticleCriteria,
) ([]model.ArticlePlus, error) {
	// See mainArticlePlusQuery
	var join string
	var joinArgs []any
	if v := criteria.FavoritedBy; v != nil {
		join = join + `
		LEFT JOIN favorites fa1 ON a.id = fa1.article_id -- at most one due to below
		LEFT JOIN users ufa1 ON ufa1.id = fa1.user_id -- AND ufa1.username = $2
		`
		//joinArgs = append(joinArgs, *v)
	}
	if v := criteria.Tag; v != nil {
		join = join + `
		LEFT JOIN article_tags at1 ON at1.article_id = a.id		
		LEFT JOIN tags t1 ON t1.id = at1.tag_id -- AND t1.name = $2
		`
		//joinArgs = append(joinArgs, *v)
	}

	var whereTuples []types.Tuple2[string, any]
	if v := criteria.FavoritedBy; v != nil {
		whereTuples = append(whereTuples, types.Tuple2[string, any]{"ufa1.username = $%d", *v})
	}
	if v := criteria.Tag; v != nil {
		whereTuples = append(whereTuples, types.Tuple2[string, any]{"t1.name = $%d", *v})
	}
	if v := criteria.Author; v != nil {
		whereTuples = append(whereTuples, types.Tuple2[string, any]{"ua.username = $%d", *v})
	}
	initialIndex := 2 + len(joinArgs)
	where, whereArgs := whereClauseFromTuples(initialIndex, whereTuples)

	additionalSql := join + where
	additionalArgs := append(joinArgs, whereArgs...)

	results, err := readArticles(ctx, tx, currUserId, additionalSql, criteria.Limit, criteria.Offset,
		additionalArgs...)
	if err != nil {
		return nil, err
	}

	return results, nil
}
