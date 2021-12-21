/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/arch/web/wgin"
)

func authenticator(pReq *http.Request) error {
	panic("todo")
}

func bodyHandlerOf[S any, T any](sfl wgin.Sfl[S, T]) gin.HandlerFunc {
	return wgin.StdFullBodyHandlerMaker[S, T](
		authenticator, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)(sfl)
}

func noBodyHandlerOf[S any, T any](sfl wgin.Sfl[S, T]) gin.HandlerFunc {
	return wgin.StdNoBodyHandlerMaker[S, T](
		authenticator, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)(sfl)
}

var articleCreateSflH = bodyHandlerOf(articleCreateSfl)
var articleDeleteSflH = noBodyHandlerOf(articleDeleteSfl)
var articleFavoriteSflH = bodyHandlerOf(articleFavoriteSfl)
var articleGetSflH = noBodyHandlerOf(articleGetSfl)
var articleUnfavoriteSflH = noBodyHandlerOf(articleUnfavoriteSfl)
var articleUpdateSflH = bodyHandlerOf(articleUpdateSfl)
var articlesFeedSflH = noBodyHandlerOf(articlesFeedSfl)
var articlesListSflH = noBodyHandlerOf(articlesListSfl)
var commentAddSflH = bodyHandlerOf(commentAddSfl)
var commentDeleteSflH = noBodyHandlerOf(commentDeleteSfl)
var commentsGetSflH = noBodyHandlerOf(commentsGetSfl)
var profileGetSflH = noBodyHandlerOf(profileGetSfl)
var tagsGetSflH = noBodyHandlerOf(tagsGetSfl)
var userAuthenticateSflH = bodyHandlerOf(userAuthenticateSfl)
var userFollowSflH = bodyHandlerOf(userFollowSfl)
var userGetCurrentSflH = noBodyHandlerOf(userGetCurrentSfl)
var userRegisterSflH = bodyHandlerOf(userRegisterSfl)
var userUnfollowSflH = noBodyHandlerOf(userUnfollowSfl)
var userUpdateSflH = bodyHandlerOf(userUpdateSfl)
