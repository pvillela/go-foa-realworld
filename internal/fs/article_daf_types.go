/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fs

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// RecCtxArticle is a type alias
//type RecCtxArticle = db.RecCtx[model.Article]

// PwArticle is a type alias
//type PwArticle = db.Pw[model.Article, RecCtxArticle]

// ArticleCreateDafT is the type of the stereotype instance for the DAF that
// creates an article.
type ArticleCreateDafT = func(ctx context.Context, article *model.Article) error

// ArticleGetBySlugDafT is the type of the stereotype instance for the DAF that
// retrieves an article by slug.
type ArticleGetBySlugDafT = func(ctx context.Context, slug string) (model.Article, error)

// ArticleUpdateDafT is the type of the stereotype instance for the DAF that
// updates an article.
type ArticleUpdateDafT = func(ctx context.Context, article *model.Article) error

// ArticleDeleteDafT is the type of the stereotype instance for the DAF that
// deletes an article.
type ArticleDeleteDafT = func(ctx context.Context, slug string) error

// ArticleGetRecentForAuthorsDafT is the type of the stereotype instance for the DAF that
// retrieves recent articles for given authors.
type ArticleGetRecentForAuthorsDafT = func(ctx context.Context, usernames []string, pLimit, pOffset *int) ([]model.Article, error)

// ArticleGetRecentFilteredDafT is the type of the stereotype instance for the DAF that
// retrieves recent articles based on filter criteria.
type ArticleGetRecentFilteredDafT = func(ctx context.Context, in rpc.ArticlesListIn) ([]model.Article, error)
