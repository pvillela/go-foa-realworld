/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package fs

import "github.com/pvillela/go-foa-realworld/internal/arch/util"

var (
	ErrDuplicateArticleSlug = util.NewErrKind("duplicate article slug \"%v\"")
	ErrArticleSlugNotFound  = util.NewErrKind("article slug \"%v\" not found")
	ErrArticleNotFound      = util.NewErrKind("article not found for Uuid %v")
	ErrCommentNotFound      = util.NewErrKind("comment not found for articleUuid %v and id %")
	ErrProfileNotFound      = util.NewErrKind("profile not found")
	ErrUserNotFound         = util.NewErrKind("user not found")
	ErrUnauthorizedUser     = util.NewErrKind("user not authorized to take this action")
	ErrAuthenticationFailed = util.NewErrKind("user authentication failed")
	ErrNotAuthenticated     = util.NewErrKind("user not authenticated")
	ErrDuplicateUser        = util.NewErrKind("user with this name already exists")
)
