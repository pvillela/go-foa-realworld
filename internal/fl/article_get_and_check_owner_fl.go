/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fl

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

// ArticleGetAndCheckOwnerFlT is the type of the stereotype instance for the flow that
// checks if a given article's author's username matches a given username.
type ArticleGetAndCheckOwnerFlT = func(
	ctx context.Context,
	tx pgx.Tx,
	slug string,
	username string,
) (model.ArticlePlus, error)

// ArticleGetAndCheckOwnerFlI implements a stereotype instance of type
// ArticleGetAndCheckOwnerFlT.
var ArticleGetAndCheckOwnerFlI ArticleGetAndCheckOwnerFlT = func(
	ctx context.Context,
	tx pgx.Tx,
	slug string,
	username string,
) (model.ArticlePlus, error) {
	userGetByNameDaf := daf.UserGetByNameExplicitTxDafI
	articleGetBySlugDaf := daf.ArticleGetBySlugDafI
	articleCheckOwnerBf := bf.ArticleCheckOwnerBfI
	return ArticleGetAndCheckOwnerFlC0(
		userGetByNameDaf,
		articleGetBySlugDaf,
		articleCheckOwnerBf,
	)(ctx, tx, slug, username)
}

// ArticleGetAndCheckOwnerFlC0 is the function that constructs a stereotype instance of type
// ArticleGetAndCheckOwnerFlT without hard-wired BF dependencies.
func ArticleGetAndCheckOwnerFlC0(
	userGetByNameDaf daf.UserGetByNameExplicitTxDafT,
	articleGetBySlugDaf daf.ArticleGetBySlugDafT,
	articleCheckOwnerBf bf.ArticleCheckOwnerBfT,
) ArticleGetAndCheckOwnerFlT {
	return func(
		ctx context.Context,
		tx pgx.Tx,
		slug string,
		username string,
	) (model.ArticlePlus, error) {
		user, _, err := userGetByNameDaf(ctx, tx, username)
		if err != nil {
			return model.ArticlePlus{}, err
		}

		article, err := articleGetBySlugDaf(ctx, tx, user.Id, slug)
		if err != nil {
			return model.ArticlePlus{}, err
		}

		if err := articleCheckOwnerBf(article, username); err != nil {
			return model.ArticlePlus{}, err
		}

		return article, err
	}
}
