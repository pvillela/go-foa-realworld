/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package bf

import "github.com/pvillela/go-foa-realworld/internal/arch/errx"

var (
	ErrMsgDuplicateArticleSlug       = "duplicate article slug \"%v\""
	ErrMsgDuplicateArticleUuid       = "duplicate article uuid %v"
	ErrMsgArticleSlugNotFound        = "article slug \"%v\" not found"
	ErrMsgArticleNotFound            = "article not found for Id %v"
	ErrMsgArticleCreateMissingFields = "article has missing fields for Create operation"
	ErrMsgArticleAlreadyFavorited    = "article with ID \"%v\" has already been favoriated"
	ErrMsgArticleWasNotFavorited     = "article with ID \"%v\" was not favorited"
	ErrMsgCommentNotFound            = "comment not found for comment ID v%"
	ErrMsgProfileNotFound            = "profile not found"
	ErrMsgTagNameAlreadyExists       = "tag name %v already exists"
	ErrMsgTagOnArticleAlreadyExists  = "tag with name %v already exists on article with slug %v"
	ErrMsgUserEmailNotFound          = "user not found for email %v"
	ErrMsgUsernameDuplicate          = "user with name \"%v\" already exists"
	ErrMsgUsernameNotFound           = "user not found for username \"%v\""
	ErrMsgUserAlreadyFollowed        = "user with username \"%v\" was already followed"
	ErrMsgUserWasNotFollowed         = "user with username \"%v\" was not followed"
	ErrMsgUnauthorizedUser           = "user \"%v\" not authorized to take this action"
	ErrMsgAuthenticationFailed       = "user authentication failed with name \"%v\" and password \"%v\""
	ErrMsgNotAuthenticated           = "user not authenticated"
	ErrMsgDuplicateUserEmail         = "user with email %v already exists"
)

var (
	ErrUnauthorizedUser = errx.NewKind("ErrUnauthorizedUser")
	ErrValidationFailed = errx.NewKind("ErrValidationFailed")
)
