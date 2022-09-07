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

// TagsGetAllDafT is the type of the stereotype instance for the DAF that
// retrieves all tags.
type TagsGetAllDafT = func(ctx context.Context, tx pgx.Tx) ([]model.Tag, error)

// TagCreateDafT is the type of the stereotype instance for the DAF that
// creates a new tag.
type TagCreateDafT = func(ctx context.Context, tx pgx.Tx, tag *model.Tag) error

// TagsAddNewDafT is the type of the stereotype instance for the DAF that
// adds a list of tags, skipping those that already exist.
type TagsAddNewDafT = func(ctx context.Context, tx pgx.Tx, names []string) error

// TagAddToArticleDafT is the type of the stereotype instance for the DAF that
// associates a tag with an article.
type TagAddToArticleDafT = func(ctx context.Context, tx pgx.Tx, tag model.Tag, article model.Article) error

// TagsAddToArticleDafT is the type of the stereotype instance for the DAF that
// associates a list of tags with an article, skipping those that are already associated.
type TagsAddToArticleDafT = func(ctx context.Context, tx pgx.Tx, names []string, article model.Article) error
