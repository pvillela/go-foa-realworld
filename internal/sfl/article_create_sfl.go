/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	rpc2 "github.com/pvillela/go-foa-realworld/rpc"

	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

// ArticleCreateSflT is the type of the stereotype instance for the service flow that
// creates an article.
type ArticleCreateSflT = func(
	ctx context.Context,
	reqCtx web.RequestContext,
	in rpc2.ArticleCreateIn,
) (rpc2.ArticleOut, error)

// ArticleCreateSflC is the function that constructs a stereotype instance of type
// ArticleCreateSflT with hard-wired stereotype dependencies.
func ArticleCreateSflC(
	cfgPvdr DefaultSflCfgPvdr,
) ArticleCreateSflT {
	return ArticleCreateSflC0(
		cfgPvdr,
		daf.UserGetByNameExplicitTxDaf,
		daf.ArticleCreateDaf,
		daf.TagsAddNewDaf,
		daf.TagsAddToArticleDaf,
	)
}

// ArticleCreateSflC0 is the function that constructs a stereotype instance of type
// ArticleCreateSflT without hard-wired stereotype dependencies.
func ArticleCreateSflC0(
	cfgPvdr DefaultSflCfgPvdr,
	userGetByNameDaf daf.UserGetByNameExplicitTxDafT,
	articleCreateDaf daf.ArticleCreateDafT,
	tagsAddNewDaf daf.TagsAddNewDafT,
	tagsAddToArticleDaf daf.TagsAddToArticleDafT,
) ArticleCreateSflT {
	db := cfgPvdr()
	return dbpgx.SflWithTransaction(db, func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		in rpc2.ArticleCreateIn,
	) (rpc2.ArticleOut, error) {
		err := in.Validate()
		if err != nil {
			return rpc2.ArticleOut{}, err
		}
		username := reqCtx.Username

		user, _, err := userGetByNameDaf(ctx, tx, username)
		if err != nil {
			return rpc2.ArticleOut{}, err
		}

		article, err := in.ToArticle(user)
		if err != nil {
			return rpc2.ArticleOut{}, err
		}

		err = articleCreateDaf(ctx, tx, &article)
		if err != nil {
			return rpc2.ArticleOut{}, err
		}

		names := article.TagList

		err = tagsAddNewDaf(ctx, tx, names)
		if err != nil {
			return rpc2.ArticleOut{}, err
		}

		err = tagsAddToArticleDaf(ctx, tx, names, article)
		if err != nil {
			return rpc2.ArticleOut{}, err
		}

		articlePlus := model.ArticlePlus_FromArticle(article, false, model.Profile_FromUser(user, false))
		articleOut := rpc2.ArticleOut_FromModel(articlePlus)

		return articleOut, nil
	})
}
