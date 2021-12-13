/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/arch/web/wgin"
)

var articleCreateSflH = wgin.PostHandlerMaker(articleCreateSfl, web.DefaultErrorHandler)
var articleDeleteSflH = wgin.PostHandlerMaker(articleDeleteSfl, web.DefaultErrorHandler)
var articleFavoriteSflH = wgin.PostHandlerMaker(articleFavoriteSfl, web.DefaultErrorHandler)
var articleGetSflH = wgin.PostHandlerMaker(articleGetSfl, web.DefaultErrorHandler)
var articleUnfavoriteSflH = wgin.PostHandlerMaker(articleUnfavoriteSfl, web.DefaultErrorHandler)
var articleUpdateSflH = wgin.PostHandlerMaker(articleUpdateSfl, web.DefaultErrorHandler)
var articlesFeedSflH = wgin.PostHandlerMaker(articlesFeedSfl, web.DefaultErrorHandler)
var articlesListSflH = wgin.PostHandlerMaker(articlesListSfl, web.DefaultErrorHandler)
var commentAddSflH = wgin.PostHandlerMaker(commentAddSfl, web.DefaultErrorHandler)
var commentDeleteSflH = wgin.PostHandlerMaker(commentDeleteSfl, web.DefaultErrorHandler)
var commentsGetSflH = wgin.PostHandlerMaker(commentsGetSfl, web.DefaultErrorHandler)
var profileGetSflH = wgin.PostHandlerMaker(profileGetSfl, web.DefaultErrorHandler)
var tagsGetSflH = wgin.PostHandlerMaker(tagsGetSfl, web.DefaultErrorHandler)
var userAuthenticateSflH = wgin.PostHandlerMaker(userAuthenticateSfl, web.DefaultErrorHandler)
var userFollowSflH = wgin.PostHandlerMaker(userFollowSfl, web.DefaultErrorHandler)
var userGet_currentSflH = wgin.PostHandlerMaker(userGetCurrentSfl, web.DefaultErrorHandler)
var userRegisterSflH = wgin.PostHandlerMaker(userRegisterSfl, web.DefaultErrorHandler)
var userUnfollowSflH = wgin.PostHandlerMaker(userUnfollowSfl, web.DefaultErrorHandler)
var userUpdateSflH = wgin.PostHandlerMaker(userUpdateSfl, web.DefaultErrorHandler)
