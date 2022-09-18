/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package boot

import (
	"github.com/pvillela/go-foa-realworld/internal/config"
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
)

///////////////////
// Config logic

var CommentsGetSflCfgAdapter = DefaultSflCfgAdapter

// CommentsGetSflBoot is the function that constructs a stereotype instance of type
// CommentsGetSflT with configuration information and hard-wired stereotype dependencies.
func CommentsGetSflBoot(
	src config.AppCfgSrc,
) sfl.CommentsGetSflT {
	return sfl.CommentsGetSflC(
		CommentsGetSflCfgAdapter(src),
		daf.CommentsGetBySlugDaf,
	)
}
