/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/fl"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// ArticleUpdateSflT is the type of the stereotype instance for the service flow that
// updates an article.
type ArticleUpdateSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	in rpc.ArticleUpdateIn,
) (rpc.ArticleOut, error)

// ArticleUpdateSflC0 is the function that constructs a stereotype instance of type
// ArticleUpdateSflT without hard-wired stereotype dependencies.
func ArticleUpdateSflC0(
	cfgSrc DefaultSflCfgSrc,
	articleGetAndCheckOwnerFl fl.ArticleGetAndCheckOwnerFlT,
	articleUpdateDaf daf.ArticleUpdateDafT,
) ArticleUpdateSflT {
	db := cfgSrc()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		in rpc.ArticleUpdateIn,
	) (rpc.ArticleOut, error) {
		err := in.Validate()
		if err != nil {
			return rpc.ArticleOut{}, err
		}

		username := reqCtx.Username
		slug := in.Article.Slug
		var zero rpc.ArticleOut

		articlePlus, _, err := articleGetAndCheckOwnerFl(ctx, tx, slug, username)
		if err != nil {
			return zero, err
		}

		article := articlePlus.ToArticle()
		updateSrc := model.ArticlePatch{
			Title:       in.Article.Title,
			Description: in.Article.Description,
			Body:        in.Article.Body,
		}
		article = article.Update(updateSrc)

		if err := articleUpdateDaf(ctx, tx, &article); err != nil {
			return rpc.ArticleOut{}, err
		}

		articlePlus = model.ArticlePlus_FromArticle(article, articlePlus.Favorited, articlePlus.Author)
		articleOut := rpc.ArticleOut_FromModel(articlePlus)

		return articleOut, err
	})
}
