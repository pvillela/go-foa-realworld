/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"github.com/pvillela/go-foa-realworld/internal/fl"
	"sync"

	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.mapdb/daf"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
)

/////
// Database-related

type txn struct{}

func (txn) Validate() error {
	return nil
}

func (txn) End() {}

func beginTxn(context string) db.Txn {
	return txn{}
}

var mapDb = mapdb.MapDb{
	Name:  "TheDatabase",
	Store: &sync.Map{},
}

/////
// DAFs and FLs

var userGetByNameDaf = daf.UserGetByEmailDafC(mapDb)
var userUpdateDaf = daf.UserUpdateDafC(mapDb)
var userFollowFl = fl.UserStartFollowingFlC(userGetByNameDaf, userUpdateDaf)
var userCreateDaf = daf.UserCreateDafC(mapDb)

var articleCreateDaf = daf.ArticleCreateDafC(mapDb)
var articleGetBySlugDaf = daf.ArticleGetBySlugDafC(mapDb)
var articleGetAndCheckOwnerFl = fl.ArticleGetAndCheckOwnerFlC(articleGetBySlugDaf)
var articleDeleteDaf = daf.ArticleDeleteDafC(mapDb)
var articleUpdateDaf = daf.ArticleUpdateDafC(mapDb)
var articleFavoriteFl = fl.ArticleFavoriteFlC(userGetByNameDaf, articleGetBySlugDaf, articleUpdateDaf)
var articleGetRecentForAuthorsDaf = daf.ArticleGetRecentForAuthorsDafC(mapDb)
var articleGetRecentFilteredDaf = daf.ArticleGetRecentFilteredDafC(mapDb)

var commentCreateDaf = daf.CommentCreateDafC(mapDb)
var commentGetByIdDaf = daf.CommentGetByIdDafC(mapDb)
var commentDeleteDaf = daf.CommentDeleteDafC(mapDb)

var tagAddDaf = daf.TagAddDafC(mapDb)
var tagGetAllDaf = daf.TagGetAllDafC(mapDb)

/////
// SFLs

var articleCreateSfl = sfl.ArticleCreateSflC(beginTxn, userGetByNameDaf, articleCreateDaf, tagAddDaf)
var articleDeleteSfl = sfl.ArticleDeleteSflC0(beginTxn, articleGetAndCheckOwnerFl, articleDeleteDaf)
var articleFavoriteSfl = sfl.ArticleFavoriteSflC0(beginTxn, articleFavoriteFl)
var articleGetSfl = sfl.ArticleGetSflC(beginTxn, userGetByNameDaf, articleGetBySlugDaf)
var articleUnfavoriteSfl = sfl.ArticleUnfavoriteSflC(beginTxn, articleFavoriteFl)
var articleUpdateSfl = sfl.ArticleUpdateSflC(beginTxn, articleGetAndCheckOwnerFl, userGetByNameDaf, articleUpdateDaf, articleGetBySlugDaf, articleCreateDaf, articleDeleteDaf)
var articlesFeedSfl = sfl.ArticlesFeedSflC(userGetByNameDaf, articleGetRecentForAuthorsDaf)
var articlesListSfl = sfl.ArticlesListSflC0(userGetByNameDaf, articleGetRecentFilteredDaf)

var commentAddSfl = sfl.CommentAddSflC0(beginTxn, userGetByNameDaf, articleGetBySlugDaf, commentCreateDaf, articleUpdateDaf)
var commentDeleteSfl = sfl.CommentDeleteSflC0(beginTxn, commentGetByIdDaf, commentDeleteDaf, articleGetBySlugDaf, articleUpdateDaf)
var commentsGetSfl = sfl.CommentsGetSflC0(articleGetBySlugDaf)

var profileGetSfl = sfl.ProfileGetSflC(userGetByNameDaf)

var tagsGetSfl = sfl.TagsGetSflC(tagGetAllDaf)

var userAuthenticateSfl = sfl.UserAuthenticateSflC(userGetByNameDaf)
var userFollowSfl = sfl.UserFollowSflC(beginTxn, userFollowFl)
var userGetCurrentSfl = sfl.UserGetCurrentSflC0(userGetByNameDaf)
var userRegisterSfl = sfl.UserRegisterSflC(beginTxn, userCreateDaf)
var userUnfollowSfl = sfl.UserUnfollowSflC(beginTxn, userFollowFl)
var userUpdateSfl = sfl.UserUpdateSflC(beginTxn, userGetByNameDaf, userUpdateDaf)
