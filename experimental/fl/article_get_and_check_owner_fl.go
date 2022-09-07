/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fl

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/experimental/bf"
	"github.com/pvillela/go-foa-realworld/experimental/model"
)

// ArticleGetAndCheckOwnerFlT is the type of the stereotype instance for the flow that
// checks if a given article's author's username matches a given username.
type ArticleGetAndCheckOwnerFlT = func(
	ctx context.Context,
	tx pgx.Tx,
	slug string,
	username string,
) (model.ArticlePlus, model.User, error)

// ArticleGetAndCheckOwnerFl implements a stereotype instance of type
// ArticleGetAndCheckOwnerFlT.
var ArticleGetAndCheckOwnerFl ArticleGetAndCheckOwnerFlT = func(
	ctx context.Context,
	tx pgx.Tx,
	slug string,
	username string,
) (model.ArticlePlus, model.User, error) {
	return ArticleGetAndCheckOwnerFlC0(
		ArticleAndUserGetFl,
		bf.ArticleCheckOwnerBf,
	)(ctx, tx, slug, username)
}

// ArticleGetAndCheckOwnerFlC0 is the function that constructs a stereotype instance of type
// ArticleGetAndCheckOwnerFlT without hard-wired BF dependencies.
func ArticleGetAndCheckOwnerFlC0(
	articleAndUserGetFl ArticleAndUserGetFlT,
	articleCheckOwnerBf bf.ArticleCheckOwnerBfT,
) ArticleGetAndCheckOwnerFlT {
	return func(
		ctx context.Context,
		tx pgx.Tx,
		slug string,
		username string,
	) (model.ArticlePlus, model.User, error) {
		article, user, err := articleAndUserGetFl(ctx, tx, slug, username)
		if err != nil {
			return model.ArticlePlus{}, model.User{}, err
		}

		if err := articleCheckOwnerBf(article, username); err != nil {
			return model.ArticlePlus{}, model.User{}, err
		}

		return article, user, err
	}
}
