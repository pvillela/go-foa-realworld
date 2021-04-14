/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package fs

import "github.com/pvillela/go-foa-realworld/internal/arch/util"

var (
	ErrDuplicateArticleSlug = util.NewErrKind("duplicate article slug \"%v\"")
	ErrDuplicateArticleUuid = util.NewErrKind("duplicate article uuid %v")
	ErrArticleSlugNotFound  = util.NewErrKind("article slug \"%v\" not found")
	ErrArticleNotFound      = util.NewErrKind("article not found for Uuid %v")
	ErrCommentNotFound      = util.NewErrKind("comment not found for articleUuid %v and id %")
	ErrProfileNotFound      = util.NewErrKind("profile not found")
	ErrUserNameNotFound     = util.NewErrKind("user not found for username \"%v\"")
	ErrUserEmailNotFound    = util.NewErrKind("user not found for email %v")
	ErrUnauthorizedUser     = util.NewErrKind("user \"%v\" not authorized to take this action")
	ErrAuthenticationFailed = util.NewErrKind("user authentication failed with name \"%v\" and password \"%v\"")
	ErrNotAuthenticated     = util.NewErrKind("user not authenticated")
	ErrDuplicateUserName    = util.NewErrKind("user with name \"%v\" already exists")
	ErrDuplicateUserEmail   = util.NewErrKind("user with email %v already exists")
)
