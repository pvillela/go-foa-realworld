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
	return wgin.BodyHandlerMaker[S, T](authenticator, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)(sfl)
}

func getHandlerOf[S any, T any](sfl wgin.Sfl[S, T], mapper func(map[string]string) (S, error)) gin.HandlerFunc {
	return wgin.GetHandlerMaker[S, T](mapper, authenticator, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)(sfl)
}

func mappedBodyHandlerOf[S any, T any](sfl wgin.Sfl[S, T], mapper func(map[string]string, *S) error) gin.HandlerFunc {
	return wgin.MappedBodyHandlerMaker[S, T](mapper, authenticator, web.DefaultReqCtxExtractor, web.DefaultErrorHandler)(sfl)
}

var articleCreateSflH = bodyHandlerOf(articleCreateSfl)
var articleDeleteSflH = wgin.DeleteHandlerMaker(articleDeleteSfl)
var articleFavoriteSflH = bodyHandlerOf(articleFavoriteSfl)
var articleGetSflH = wgin.GetHandlerMaker(articleGetSfl)
var articleUnfavoriteSflH = wgin.DeleteHandlerMaker(articleUnfavoriteSfl)
var articleUpdateSflH = wgin.PutHandlerMaker(articleUpdateSfl)
var articlesFeedSflH = wgin.GetHandlerMaker(articlesFeedSfl)
var articlesListSflH = wgin.GetHandlerMaker(articlesListSfl)
var commentAddSflH = bodyHandlerOf(commentAddSfl)
var commentDeleteSflH = wgin.DeleteHandlerMaker(commentDeleteSfl)
var commentsGetSflH = wgin.GetHandlerMaker(commentsGetSfl)
var profileGetSflH = wgin.GetHandlerMaker(profileGetSfl)
var tagsGetSflH = wgin.GetHandlerMaker(tagsGetSfl)
var userAuthenticateSflH = bodyHandlerOf(userAuthenticateSfl)
var userFollowSflH = bodyHandlerOf(userFollowSfl)
var userGetCurrentSflH = wgin.GetHandlerMaker(userGetCurrentSfl)
var userRegisterSflH = bodyHandlerOf(userRegisterSfl)
var userUnfollowSflH = wgin.DeleteHandlerMaker(userUnfollowSfl)
var userUpdateSflH = wgin.PutHandlerMaker(userUpdateSfl)
