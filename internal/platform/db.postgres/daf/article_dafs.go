/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
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
	user model.User,
	slug string,
) (model.ArticlePlus, error) {
	selectJoin := `
	SELECT a.*, u.username, u.bio, u.image, f1.user_id as favorited, fo.follower_id as following, t.name 
		FROM articles a
	LEFT JOIN users u ON a.author_id = u.id
	LEFT JOIN favorites f1 ON a.id = f1.article_id AND $1 = f1.user_id -- no product effect
	LEFT JOIN favorites f2 ON a.id = f2.article_id AND $2 = f2.user_id -- no product effect
	LEFT JOIN followings fo ON fo.follower_id = $1 AND a.author_id = fo.followee_id -- no product effect
	LEFT JOIN article_tags at ON a.id = at.article_id -- product effect
	LEFT JOIN tags t ON at.tag_id = t.id
	`
	where := `
	WHERE slug = $2 -- AND f2.user_id IS NOT NULL
	`
	orderBy := `
	ORDER BY a.created_at DESC, t.name
	`
	sql := selectJoin + where + orderBy
	args := []any{user.Id, user.Id, slug}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return model.ArticlePlus{}, errx.ErrxOf(err)
	}
	defer rows.Close()

	tags := []string{}
	var record struct {
		article   model.Article `db:""`
		profile   model.Profile `db:""`
		following *uint
		favorited *uint
		tag       string `db:"t.name"`
	}
	for rows.Next() {
		err = pgxscan.ScanRow(&record, rows)
		if err != nil {
			return model.ArticlePlus{}, errx.ErrxOf(err)
		}
		tags = append(tags, record.tag)
	}

	ra := &record.article
	rp := &record.profile
	if record.following != nil {
		rp.Following = true
	}
	favorited := false
	if record.favorited != nil {
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

	return articlePlus, nil
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
var ArticlesFeedDaf0 ArticlesFeedDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	user model.User,
	optLimit *int,
	optOffset *int,
) ([]model.Article, error) {
	mainSql := `
	SELECT a.* from articles a
	INNER JOIN followings f ON $1 = f.follower_id AND a.author_id = f.followee_id
	ORDER BY created_at DESC
	`
	args := []any{user.Id}

	return readArticles(ctx, tx, mainSql, optLimit, optOffset, args...)
}

// SQL WHERE clauses used by ArticlesListDaf
var clauses = map[string]string{
	"tag": `
		id IN (select article_id from article_tags where tag_id in (
			select id from tags where name = $%d)
		)`,
	"author": `
		author_id = (select id from users where username = $%d)`,
	"favoritedBy": `
		id IN (select article_id from favorites where user_id = (
			select id from users where username = $%d)
		)`,
}

// ArticlesListDafT implements a stereotype instance of type
// ArticlesListDafTT.
var ArticlesListDaf ArticlesListDafT = func(
	ctx context.Context,
	tx pgx.Tx,
	criteria model.ArticleCriteria,
) ([]model.Article, error) {
	namedArgs := make([]util.Tuple2[string, any], 0)
	if v := criteria.Tag; v != nil {
		namedArgs = append(namedArgs, util.Tuple2[string, any]{"tag", *v})
	}
	if v := criteria.Author; v != nil {
		namedArgs = append(namedArgs, util.Tuple2[string, any]{"aurhor", *v})
	}
	if v := criteria.FavoritedBy; v != nil {
		namedArgs = append(namedArgs, util.Tuple2[string, any]{"favoritedBy", *v})
	}

	wheres := make([]string, len(namedArgs))
	args := make([]any, len(namedArgs))
	for i, nv := range namedArgs {
		wheres[i] = fmt.Sprintf(clauses[nv.X1], i+1)
		args[i] = nv.X2
	}
	where := strings.Join(wheres, " AND ")
	if len(wheres) != 0 {
		where = " WHERE " + where
	}

	mainSql := "SELECT * from articles" + where + " ORDER BY created_at DESC"

	return readArticles(ctx, tx, mainSql, criteria.Limit, criteria.Offset, args...)
}

func readArticles(
	ctx context.Context,
	tx pgx.Tx,
	mainSql string,
	optLimit *int,
	optOffset *int,
	args ...any,
) ([]model.Article, error) {
	limit := limitDefault
	if optLimit != nil {
		limit = *optLimit
	}
	offset := offsetDefault
	if optOffset != nil {
		offset = *optOffset
	}
	return dbpgx.ReadMany[model.Article](ctx, tx, mainSql, limit, offset, args...)
}
