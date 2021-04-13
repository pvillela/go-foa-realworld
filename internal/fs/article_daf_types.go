/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// PwArticle is a wrapper of the model.User entity
// containing context information required for persistence purposes.
type PwArticle struct {
	db.RecCtx
	Entity model.Article
}

type ArticleCreateDafT = func(article model.Article, txn mapdb.Txn) (db.RecCtx, error)

type ArticleGetBySlugDafT = func(slug string) (model.Article, db.RecCtx, error)

type ArticleUpdateDafT = func(article model.Article, recCtx db.RecCtx, txn mapdb.Txn) (db.RecCtx, error)

type ArticleDeleteDafT = func(slug string, txn mapdb.Txn) error

type ArticleGetByAuthorsOrderedByMostRecentDafT = func(usernames []string) ([]model.Article, error)

type ArticleGetRecentFilteredDafT = func(filters []model.ArticleFilter) ([]model.Article, error)
