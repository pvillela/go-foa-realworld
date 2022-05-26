/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/stretchr/testify/assert"
	"testing"
)

//import (
//	"errors"
//	"fmt"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)

func TestArticleCreateDafI(t *testing.T) {
	dafTester(func(t *testing.T, ctx context.Context, tx pgx.Tx) {
		currUserId := users[0].Id

		{
			author := users[1]

			criteria := model.ArticleCriteria{
				Tag:         nil,
				Author:      &author.Username,
				FavoritedBy: nil,
				Limit:       nil,
				Offset:      nil,
			}
			articlePluses, err := daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
			errx.PanicOnError(err)

			coreInfoReturned := util.SliceMap(articlePluses, func(_ int, ap model.ArticlePlus) model.ArticlePlus {
				return model.ArticlePlus{
					Slug:        ap.Slug,
					Author:      ap.Author,
					Title:       ap.Title,
					Description: ap.Description,
					Body:        ap.Body,
					TagList:     ap.TagList,
				}
			})

			coreInfoExpected := util.SliceMap(articles, func(_ int, a model.Article) model.ArticlePlus {
				return model.ArticlePlus{
					Slug:        a.Slug,
					Author:      model.Profile_FromUser(&author, false),
					Title:       a.Title,
					Description: a.Description,
					Body:        a.Body,
					TagList:     a.TagList,
				}
			})

			//fmt.Println("\ncoreInfoReturned:", coreInfoReturned)
			//fmt.Println("\ncoreInfoExpected:", coreInfoExpected)

			assert.ElementsMatch(t, coreInfoExpected, coreInfoReturned)
		}

		{
			criteria := model.ArticleCriteria{
				Tag:         util.PointerFromValue("FOOTAG"),
				Author:      nil,
				FavoritedBy: nil,
				Limit:       nil,
				Offset:      nil,
			}
			articlePluses, err := daf.ArticlesListDafI(ctx, tx, currUserId, criteria)
			errx.PanicOnError(err)

			coreInfoReturned := util.SliceMap(articlePluses, func(_ int, ap model.ArticlePlus) model.ArticlePlus {
				return model.ArticlePlus{
					Slug:        ap.Slug,
					Author:      ap.Author,
					Title:       ap.Title,
					Description: ap.Description,
					Body:        ap.Body,
					TagList:     ap.TagList,
				}
			})

			var coreInfoExpected []model.ArticlePlus

			//fmt.Println("\narticlesListDaf - by tag:", articlePluses)

			assert.ElementsMatch(t, coreInfoExpected, coreInfoReturned)
		}

		{
			articleFromDb, err := daf.ArticleGetBySlugDafI(ctx, tx, currUserId, articles[1].Slug)
			errx.PanicOnError(err)
			fmt.Println("\nArticleGetBySlugDaf:", articleFromDb)
		}

		{
			pArticle := &articles[0]
			pArticle.Title = "A very interesting subject"
			err := daf.ArticleUpdateDafI(ctx, tx, pArticle)
			errx.PanicOnError(err)
			fmt.Println("ArticleUpdateDaf:", pArticle)
		}

		{
			articlePluses, err := daf.ArticlesFeedDafI(ctx, tx, currUserId, nil, nil)
			errx.PanicOnError(err)
			fmt.Println("\nArticlesFeedDaf:", articlePluses)
		}
	})(t)
}
