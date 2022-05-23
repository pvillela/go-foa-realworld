/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// RecCtxArticle is a type alias
//type RecCtxArticle = db.RecCtx[model.Article]

// PwArticle is a type alias
//type PwArticle = db.Pw[model.Article, RecCtxArticle]

// ArticleCreateDafT is the type of the stereotype instance for the DAF that
// creates an article.
type ArticleCreateDafT = func(ctx context.Context, tx pgx.Tx, article *model.Article) error

// ArticleGetBySlugDafT is the type of the stereotype instance for the DAF that
// retrieves an article by slug.
type ArticleGetBySlugDafT = func(ctx context.Context, tx pgx.Tx, currUserId uint, slug string) (model.ArticlePlus, error)

// ArticleUpdateDafT is the type of the stereotype instance for the DAF that
// updates an article.
type ArticleUpdateDafT = func(ctx context.Context, tx pgx.Tx, article *model.Article) error

// ArticleDeleteDafT is the type of the stereotype instance for the DAF that
// deletes an article.
type ArticleDeleteDafT = func(ctx context.Context, tx pgx.Tx, slug string) error

// ArticlesFeedDafT is the type of the stereotype instance for the DAF that
// queries for all articles authored by other users followed by
// the current user, with optional limit and offset pagination parameters.
type ArticlesFeedDafT = func(ctx context.Context, tx pgx.Tx, currUserId uint, optLimit *int, optOffset *int) ([]model.ArticlePlus, error)

// ArticlesListDafT is the type of the stereotype instance for the DAF that
// retrieve recent articles based on a set of query parameters.
type ArticlesListDafT = func(ctx context.Context, tx pgx.Tx, currUserId uint, criteria model.ArticleCriteria) ([]model.ArticlePlus, error)
