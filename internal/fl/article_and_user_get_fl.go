/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fl

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

// ArticleAndUserGetFlT is the type of the stereotype instance for the flow that
// retrieves an article by slug and a user by username. It assumes that the username
// is that of the current user.
type ArticleAndUserGetFlT = func(
	ctx context.Context,
	tx pgx.Tx,
	slug string,
	username string,
) (model.ArticlePlus, model.User, error)

// ArticleAndUserGetFl implements a stereotype instance of type
// ArticleAndUserGetFlT.
var ArticleAndUserGetFl ArticleAndUserGetFlT = func(
	ctx context.Context,
	tx pgx.Tx,
	slug string,
	username string,
) (model.ArticlePlus, model.User, error) {
	userGetByNameDaf := daf.UserGetByNameExplicitTxDaf
	articleGetBySlugDaf := daf.ArticleGetBySlugDaf
	return ArticleAndUserGetFlC0(
		userGetByNameDaf,
		articleGetBySlugDaf,
	)(ctx, tx, slug, username)
}

// ArticleAndUserGetFlC0 is the function that constructs a stereotype instance of type
// ArticleAndUserGetFlT without hard-wired BF dependencies.
func ArticleAndUserGetFlC0(
	userGetByNameDaf daf.UserGetByNameExplicitTxDafT,
	articleGetBySlugDaf daf.ArticleGetBySlugDafT,
) ArticleAndUserGetFlT {
	return func(
		ctx context.Context,
		tx pgx.Tx,
		slug string,
		username string,
	) (model.ArticlePlus, model.User, error) {
		user, _, err := userGetByNameDaf(ctx, tx, username)
		if err != nil {
			return model.ArticlePlus{}, model.User{}, err
		}

		article, err := articleGetBySlugDaf(ctx, tx, user.Id, slug)
		if err != nil {
			return model.ArticlePlus{}, model.User{}, err
		}

		return article, user, nil
	}
}
