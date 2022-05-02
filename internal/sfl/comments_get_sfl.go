/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// CommentsGetSflC is the type of the stereotype instance for the service flow that
// retrieves the comments of an article.
type CommentsGetSflT = func(ctx context.Context, slug string) (rpc.CommentsOut, error)

// CommentsGetSflC is the function that constructs a stereotype instance of type
// CommentsGetSflT.
func CommentsGetSflC(
	articleGetBySlugDaf daf.ArticleGetBySlugDafT,
) CommentsGetSflT {
	return func(_ context.Context, slug string) (rpc.CommentsOut, error) {
		var zero rpc.CommentsOut

		article, _, err := articleGetBySlugDaf(slug)
		if err != nil {
			return zero, err
		}

		if article.Comments == nil {
			article.Comments = []model.Comment{}
		}

		commentsOut := rpc.CommentsOut_FromModel(article.Comments)

		return commentsOut, nil
	}
}
