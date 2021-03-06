/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package fs

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// PwArticle is a wrapper of the model.User entity
// containing context information required for persistence purposes.
type PwArticle struct {
	db.RecCtx
	Entity model.Article
}

type ArticleCreateDafT = func(article model.Article, txn db.Txn) (db.RecCtx, error)

type ArticleGetBySlugDafT = func(slug string) (model.Article, db.RecCtx, error)

type ArticleUpdateDafT = func(article model.Article, recCtx db.RecCtx, txn db.Txn) (db.RecCtx, error)

type ArticleDeleteDafT = func(slug string, txn db.Txn) error

type ArticleGetRecentForAuthorsDafT = func(usernames []string, pLimit, pOffset *int) ([]model.Article, error)

type ArticleGetRecentFilteredDafT = func(in rpc.ArticlesListIn) ([]model.Article, error)
