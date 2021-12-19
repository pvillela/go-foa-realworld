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

var articleCreateSflH = wgin.PostHandlerMaker(articleCreateSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var articleDeleteSflH = wgin.DeleteHandlerMaker(articleDeleteSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var articleFavoriteSflH = wgin.PostHandlerMaker(articleFavoriteSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var articleGetSflH = wgin.GetHandlerMaker(articleGetSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var articleUnfavoriteSflH = wgin.DeleteHandlerMaker(articleUnfavoriteSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var articleUpdateSflH = wgin.PutHandlerMaker(articleUpdateSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var articlesFeedSflH = wgin.GetHandlerMaker(articlesFeedSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var articlesListSflH = wgin.GetHandlerMaker(articlesListSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var commentAddSflH = wgin.PostHandlerMaker(commentAddSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var commentDeleteSflH = wgin.DeleteHandlerMaker(commentDeleteSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var commentsGetSflH = wgin.GetHandlerMaker(commentsGetSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var profileGetSflH = wgin.GetHandlerMaker(profileGetSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var tagsGetSflH = wgin.GetHandlerMaker(tagsGetSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var userAuthenticateSflH = wgin.PostHandlerMaker(userAuthenticateSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var userFollowSflH = wgin.PostHandlerMaker(userFollowSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var userGetCurrentSflH = wgin.GetHandlerMaker(userGetCurrentSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var userRegisterSflH = wgin.PostHandlerMaker(userRegisterSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var userUnfollowSflH = wgin.DeleteHandlerMaker(userUnfollowSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
var userUpdateSflH = wgin.PutHandlerMaker(userUpdateSfl, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)
