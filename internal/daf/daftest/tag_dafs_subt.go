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
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx/dbpgxtest"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/rpc"
	"github.com/stretchr/testify/assert"
	"testing"
)

var tagDafsSubt = dbpgxtest.TestWithTransaction(func(ctx context.Context, tx pgx.Tx, t *testing.T) {

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
		err := daf.TagCreateDaf(ctx, tx, &tag)
		assert.NoError(t, err)

		mdb.TagUpsert(name, tag)
	}

	for _, taggedSlug := range taggedSlugs {
		name := taggedSlug.tagName
		slug := taggedSlug.slug
		tag := mdb.TagGet(name)
		article := mdb.ArticleGetBySlug(slug)

		err := daf.TagAddToArticleDaf(ctx, tx, tag, article)
		assert.NoError(t, err)

		mdb.TagAssignToSlug(name, slug)
	}

	// Tests

	{
		msg := "TagCreateDaf - attempt to reinsert existing tag"

		tagName := tagName1

		tag := mdb.TagGet(tagName)

		// --> start nested transaction to avoid invalidating main transaction tx
		subTx, err := tx.Begin(ctx)
		assert.NoError(t, err)

		err = daf.TagCreateDaf(ctx, tx, &tag)

		retErrxKind := dbpgx.ClassifyError(err)
		expErrxKind := dbpgx.DbErrUniqueViolation

		err = subTx.Rollback(ctx)
		assert.NoError(t, err)
		// <-- rolled back nested transaction

		assert.Equal(t, expErrxKind, retErrxKind, msg)
	}

	{
		msg := "TagsGetAllDaf"

		returned, err := daf.TagsGetAllDaf(ctx, tx)
		assert.NoError(t, err)
		//fmt.Println("Tags from database:", returned, "\n")

		expected := mdb.TagGetAll()

		assert.ElementsMatchf(t, expected, returned, msg)
	}

	{
		msg := "all articles are properly tagged"

		currUsername := username1

		currUserId := mdb.UserGetByName(currUsername).Id

		criteria := rpc.ArticleCriteria{
			Tag:         nil,
			Author:      nil,
			FavoritedBy: nil,
			Limit:       nil,
			Offset:      nil,
		}
		returnedArticlePluses, err := daf.ArticlesListDaf(ctx, tx, currUserId, criteria)
		assert.NoError(t, err)

		returnedArticleTagsMap := articlePlusesToArticleTagsMap(returnedArticlePluses)
		expectedArticleTagsMap := articlePlusesToArticleTagsMap(mdb.ArticlePlusGetAll(currUsername))

		assert.Equal(t, expectedArticleTagsMap, returnedArticleTagsMap, msg)
	}

	{
		msg := "ArticlesListDaf by tag"

		currUsername := username1

		currUserId := mdb.UserGetByName(currUsername).Id

		criteria := rpc.ArticleCriteria{
			Tag:         util.PointerFromValue(tagName1),
			Author:      nil,
			FavoritedBy: nil,
			Limit:       nil,
			Offset:      nil,
		}
		returnedArticlePluses, err := daf.ArticlesListDaf(ctx, tx, currUserId, criteria)
		assert.NoError(t, err)

		returnedArticleTagsMap := articlePlusesToArticleTagsMap(returnedArticlePluses)
		expectedArticlePluses := util.SliceFilter(mdb.ArticlePlusGetAll(currUsername), func(ap model.ArticlePlus) bool {
			return util.SliceContains(ap.TagList, tagName1)
		})
		expectedArticleTagsMap := articlePlusesToArticleTagsMap(expectedArticlePluses)

		assert.Equal(t, expectedArticleTagsMap, returnedArticleTagsMap, msg)
	}

	{
		msg := "TagsAddNewDaf - add some more tags to database"

		newNames := []string{"ZZZ", "WWW"}
		names := append(newNames, tagName1)

		{
			err := daf.TagsAddNewDaf(ctx, tx, names)
			assert.NoError(t, err)

			returned, err := daf.TagsGetAllDaf(ctx, tx)
			assert.NoError(t, err)

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

			err := daf.TagsAddToArticleDaf(ctx, tx, names, article)
			assert.NoError(t, err)

			for _, name := range names {
				mdb.TagAssignToSlug(name, article.Slug)
			}

			{
				msg := "all articles are properly tagged again"

				currUsername := username1

				currUserId := mdb.UserGetByName(currUsername).Id

				criteria := rpc.ArticleCriteria{
					Tag:         nil,
					Author:      nil,
					FavoritedBy: nil,
					Limit:       nil,
					Offset:      nil,
				}
				returnedArticlePluses, err := daf.ArticlesListDaf(ctx, tx, currUserId, criteria)
				assert.NoError(t, err)

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
