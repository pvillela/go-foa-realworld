/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import "github.com/gin-gonic/gin"

func SetRoutes(r gin.IRouter) {
	r.POST("/api/users/login", userAuthenticateSflH)                // Authentication
	r.POST("/api/users", userRegisterSflH)                          // Registration
	r.GET("/api/user", userGetCurrentSflH)                          // Get current user
	r.PUT("/api/user", userUpdateSflH)                              // Update user
	r.GET("/api/profiles/:username", profileGetSflH)                // Get profile
	r.POST("/api/profiles/:username/follow", userFollowSflH)        // Follow user
	r.DELETE("/api/profiles/:username/follow", userUnfollowSflH)    // Unfollow user
	r.GET("/api/articles", articlesListSflH)                        // List articlces
	r.GET("/api/articles/feed", articlesFeedSflH)                   // Feed articles
	r.GET("/api/articles/:slug", articleGetSflH)                    // Get article
	r.POST("/api/articles", articleCreateSflH)                      // Create article
	r.PUT("/api/articles/:slug", articleUpdateSflH)                 // Update article
	r.DELETE("/api/articles/:slug", articleDeleteSflH)              // Delete article
	r.POST("/api/articles/:slug/comments", commentAddSflH)          // Add comments to an article
	r.GET("/api/articles/:slug/comments", commentsGetSflH)          // Get comments from an article
	r.DELETE("/api/articles/:slug/comments/:id", commentDeleteSflH) // Delete comment
	r.POST("/api/articles/:slug/favorite", articleFavoriteSflH)     // Favorite article
	r.DELETE("/api/articles/:slug/favorite", articleUnfavoriteSflH) // Unfavorite article
	r.GET("/api/tags", tagsGetSflH)                                 // Get tags
}
