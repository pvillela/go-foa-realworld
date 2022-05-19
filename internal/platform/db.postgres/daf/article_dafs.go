/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/model"
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
	// See mainArticlePlusQuery
	whereTuples := []util.Tuple2[string, any]{
		{"a.slug = $%d", slug},
	}
	where, whereArgs := whereClauseFromTuples(2, whereTuples)

	results, err := readArticles(ctx, tx, currUserId, where, nil, nil, whereArgs...)
	if err != nil {
		return model.ArticlePlus{}, err
	}

	if len(results) == 0 {
		return model.ArticlePlus{}, bf.ErrArticleSlugNotFound.Make(nil, slug)
	}
	if len(results) > 1 {
		util.PanicOnError(errx.NewErrx(nil,
			fmt.Sprintf("Found multiple articles with same slug '%v'", slug)))
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
	// See mainArticlePlusQuery
	whereTuples := []util.Tuple2[string, any]{
		{"fo.follower_id = $%d", currUserId},
	}
	where, whereArgs := whereClauseFromTuples(2, whereTuples)

	results, err := readArticles(ctx, tx, currUserId, where, optLimit, optOffset, whereArgs...)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// ArticlesListDafT implements a stereotype instance of type
// ArticlesListDafTT.
var ArticlesListDaf ArticlesListDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	currUserId uint,
	criteria model.ArticleCriteria,
) ([]model.ArticlePlus, error) {
	// See mainArticlePlusQuery
	var join string
	var joinArgs []any
	if v := criteria.FavoritedBy; v != nil {
		join = `
		LEFT JOIN favorites fa1 ON a.id = fa1.article_id -- at most one due to below
		LEFT JOIN users ufa1 ON ufa1.id = fa1.user_id AND ufa1.username = $2
		`
		joinArgs = []any{*v}
	}

	var whereTuples []util.Tuple2[string, any]
	if v := criteria.FavoritedBy; v != nil {
		whereTuples = append(whereTuples, util.Tuple2[string, any]{"ufa1.username = $%d", *v})
	}
	if v := criteria.Tag; v != nil {
		whereTuples = append(whereTuples, util.Tuple2[string, any]{"t.name = $%d", *v})
	}
	if v := criteria.Author; v != nil {
		whereTuples = append(whereTuples, util.Tuple2[string, any]{"ua.username = $%d", *v})
	}
	var initialIndex int
	if len(join) == 0 {
		initialIndex = 2
	} else {
		initialIndex = 3
	}
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
