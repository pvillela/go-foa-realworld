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
	"time"
)

// ArticleCreateDafI implements a stereotype instance of type
// ArticleCreateDafT.
var ArticleCreateDafI ArticleCreateDafT = func(
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

// ArticleGetBySlugDafI implements a stereotype instance of type
// ArticleGetBySlugDafT.
var ArticleGetBySlugDafI ArticleGetBySlugDafT = func(
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

// ArticleUpdateDafI implements a stereotype instance of type
// ArticleUpdateDafT.
var ArticleUpdateDafI ArticleUpdateDafT = func(
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

// ArticleAdjustFavoritesCountDafI implements a stereotype instance of type
// ArticleAdjustFavoritesCountDafT.
var ArticleAdjustFavoritesCountDafI ArticleAdjustFavoritesCountDafT = func(
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

// ArticleDeleteDafI implements a stereotype instance of type
// ArticleDeleteDafT.
var ArticleDeleteDafI ArticleDeleteDafT = func(
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

// ArticlesFeedDafI implements a stereotype instance of type
// ArticlesFeedDafT.
var ArticlesFeedDafI ArticlesFeedDafT = func(
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

// ArticlesListDafI implements a stereotype instance of type
// ArticlesListDafT.
var ArticlesListDafI ArticlesListDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	currUserId uint,
	criteria model.ArticleCriteria,
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
