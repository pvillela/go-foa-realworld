/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package bf

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
)

var (
	ErrDuplicateArticleSlug       = errx.NewKind("duplicate article slug \"%v\"")
	ErrDuplicateArticleUuid       = errx.NewKind("duplicate article uuid %v")
	ErrArticleSlugNotFound        = errx.NewKind("article slug \"%v\" not found")
	ErrArticleNotFound            = errx.NewKind("article not found for Id %v")
	ErrArticleCreateMissingFields = errx.NewKind("article has missing fields for Create operation")
	ErrArticleUpdateMissingFields = errx.NewKind("article has missing fields for Update operation")
	ErrCommentNotFound            = errx.NewKind("comment not found for articleUuid %v and id %")
	ErrProfileNotFound            = errx.NewKind("profile not found")
	ErrUserNameNotFound           = errx.NewKind("user not found for username \"%v\"")
	ErrUserEmailNotFound          = errx.NewKind("user not found for email %v")
	ErrUnauthorizedUser           = errx.NewKind("user \"%v\" not authorized to take this action")
	ErrAuthenticationFailed       = errx.NewKind("user authentication failed with name \"%v\" and password \"%v\"")
	ErrNotAuthenticated           = errx.NewKind("user not authenticated")
	ErrDuplicateUserName          = errx.NewKind("user with name \"%v\" already exists")
	ErrDuplicateUserEmail         = errx.NewKind("user with email %v already exists")
)
