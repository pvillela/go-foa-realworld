/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tagDafsSubt = dbpgx.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {

	// Create tags.

	const (
		tagName1 = "FOOTAG"
		tagName2 = "BARTAG"
	)

	tagNames := []string{
		tagName1,
		tagName2,
	}

	type taggedSlugT struct {
		tagName string
		slug    string
	}

	taggedSlugs := []taggedSlugT{
		{
			tagName: tagName1,
			slug:    slug2,
		},
		{
			tagName: tagName2,
			slug:    slug2,
		},
	}

	for _, name := range tagNames {
		tag := model.Tag{
			Name: name,
		}
		err := daf.TagCreateDafI(ctx, tx, &tag)
		errx.PanicOnError(err)

		mdb.TagUpsert(name, tag)
	}

	for _, taggedSlug := range taggedSlugs {
		name := taggedSlug.tagName
		slug := taggedSlug.slug
		tag := mdb.TagGet(name)
		article := mdb.ArticleGetBySlug(slug)

		err := daf.TagAddToArticleDafI(ctx, tx, tag, article)
		errx.PanicOnError(err)

		mdb.TagAssignToSlug(name, slug)
	}

	// Tests

	{
		msg := "TagCreateDafI - attempt to reinsert existing tag"

		tagName := tagName1

		tag := mdb.TagGet(tagName)

		// --> start nested transaction to avoid invalidating main transaction tx
		subTx, err := tx.Begin(ctx)
		errx.PanicOnError(err)

		err = daf.TagCreateDafI(ctx, tx, &tag)

		retErrxKind := dbpgx.ClassifyError(err)
		expErrxKind := dbpgx.DbErrUniqueViolation

		err = subTx.Rollback(ctx)
		errx.PanicOnError(err)
		// <-- rolled back nested transaction

		assert.Equal(t, expErrxKind, retErrxKind, msg)
	}

	{
		msg := "TagsGetAllDafI"

		returned, err := daf.TagsGetAllDafI(ctx, tx)
		errx.PanicOnError(err)
		//fmt.Println("Tags from database:", returned, "\n")

		expected := mdb.TagGetAll()

		assert.ElementsMatchf(t, expected, returned, msg)
	}

	{
		msg := "all articles are properly tagged"

		currUsername := username1

		currUserId := mdb.UserGetByName(currUsername).Id

		criteria := model.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: nil,
			Limit:       nil,
			Offset:      nil,
		}
		returnedArticlePluses, err := daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
		errx.PanicOnError(err)

		returnedArticleTagsMap := articlePlusesToArticleTagsMap(returnedArticlePluses)
		expectedArticleTagsMap := articlePlusesToArticleTagsMap(mdb.ArticlePlusGetAll(currUsername))

		assert.Equal(t, expectedArticleTagsMap, returnedArticleTagsMap, msg)
	}

	{
		msg := "ArticlesListDafI by tag"

		currUsername := username1

		currUserId := mdb.UserGetByName(currUsername).Id

		criteria := model.ArticleCriteria{
			Tag:         util.PointerFromValue(tagName1),
			Author:      nil,
			FavoritedBy: nil,
			Limit:       nil,
			Offset:      nil,
		}
		returnedArticlePluses, err := daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
		errx.PanicOnError(err)

		returnedArticleTagsMap := articlePlusesToArticleTagsMap(returnedArticlePluses)
		expectedArticlePluses := util.SliceFilter(mdb.ArticlePlusGetAll(currUsername), func(ap model.ArticlePlus) bool {
			return util.SliceContains(ap.TagList, tagName1)
		})
		expectedArticleTagsMap := articlePlusesToArticleTagsMap(expectedArticlePluses)

		assert.Equal(t, expectedArticleTagsMap, returnedArticleTagsMap, msg)
	}

	{
		msg := "TagsAddNewDafI - add some more tags to database"

		newNames := []string{"ZZZ", "WWW"}
		names := append(newNames, tagName1)

		{
			err := daf.TagsAddNewDafI(ctx, tx, names)
			errx.PanicOnError(err)

			returned, err := daf.TagsGetAllDafI(ctx, tx)
			errx.PanicOnError(err)

			originalNames := mdb.TagGetAllNames()

			for _, tag := range returned {
				mdb.TagUpsert(tag.Name, tag)
			}

			actual := util.SliceMap(returned, func(tag model.Tag) string {
				return tag.Name
			})

			expected := append(originalNames, newNames...)

			assert.ElementsMatchf(t, expected, actual, msg)
		}

		// Add the new tags to an article.
		{
			taggedSlug := slug2

			article := mdb.ArticleGetBySlug(taggedSlug)

			err := daf.TagsAddToArticleDafI(ctx, tx, names, article)
			errx.PanicOnError(err)

			for _, name := range names {
				mdb.TagAssignToSlug(name, article.Slug)
			}

			{
				msg := "all articles are properly tagged again"

				currUsername := username1

				currUserId := mdb.UserGetByName(currUsername).Id

				criteria := model.ArticleCriteria{
					Tag:         nil,
					Author:      nil,
					FavoritedBy: nil,
					Limit:       nil,
					Offset:      nil,
				}
				returnedArticlePluses, err := daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
				errx.PanicOnError(err)

				returnedArticleTagsMap := articlePlusesToArticleTagsMap(returnedArticlePluses)
				expectedArticleTagsMap := articlePlusesToArticleTagsMap(mdb.ArticlePlusGetAll(currUsername))

				assert.Equal(t, expectedArticleTagsMap, returnedArticleTagsMap, msg)
			}
		}
	}
})

// articlePlusesToArticleTagsMap transforms a list of ArticlePlus into a map from
// slug to the set of tag names for the article.
func articlePlusesToArticleTagsMap(aps []model.ArticlePlus) map[string]map[string]bool {
	articleTagsMap := map[string]map[string]bool{}
	for _, ap := range aps {
		tagSet := map[string]bool{}
		for _, name := range ap.TagList {
			tagSet[name] = true
		}
		articleTagsMap[ap.Slug] = tagSet
	}
	return articleTagsMap
}
