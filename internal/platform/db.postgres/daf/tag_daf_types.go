/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/model"
)

// TagGetAllDafT is the type of the stereotype instance for the DAF that
// retrieves all tags.
type TagGetAllDafT = func(ctx context.Context, tx pgx.Tx) ([]string, error)

// TagCreateDafT is the type of the stereotype instance for the DAF that
// creates a new tag.
type TagCreateDafT = func(ctx context.Context, tx pgx.Tx, tag *model.Tag) error

// TagAddToArticleDafT is the type of the stereotype instance for the DAF that
// associates a tag with an article.
type TagAddToArticleDafT = func(ctx context.Context, tx pgx.Tx, tag model.Tag, article model.Article) error
