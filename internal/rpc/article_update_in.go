/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

type ArticleUpdateIn struct {
	Article struct {
		Slug        string
		Title       *string // optional
		Description *string // optional
		Body        *string // optional
	}
}
