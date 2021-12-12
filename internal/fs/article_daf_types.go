/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// PwArticle is a type alias
type PwArticle = db.Pw[model.Article]

// RecCtxArticle is a type alias
type RecCtxArticle = db.RecCtx[model.Article]

// ArticleCreateDafT is the type of the stereotype instance for the DAF that
// creates an article.
type ArticleCreateDafT = func(article model.Article, txn db.Txn) (RecCtxArticle, error)

// ArticleGetBySlugDafT is the type of the stereotype instance for the DAF that
// retrieves an article by slug.
type ArticleGetBySlugDafT = func(slug string) (model.Article, RecCtxArticle, error)

// ArticleUpdateDafT is the type of the stereotype instance for the DAF that
// updates an article.
type ArticleUpdateDafT = func(article model.Article, recCtx RecCtxArticle, txn db.Txn) (RecCtxArticle, error)

// ArticleDeleteDafT is the type of the stereotype instance for the DAF that
// deletes an article.
type ArticleDeleteDafT = func(slug string, txn db.Txn) error

// ArticleGetRecentForAuthorsDafT is the type of the stereotype instance for the DAF that
// retrieves recent articles for given authors.
type ArticleGetRecentForAuthorsDafT = func(usernames []string, pLimit, pOffset *int) ([]model.Article, error)

// ArticleGetRecentFilteredDafT is the type of the stereotype instance for the DAF that
// retrieves recent articles based on filter criteria.
type ArticleGetRecentFilteredDafT = func(in rpc.ArticlesListIn) ([]model.Article, error)
