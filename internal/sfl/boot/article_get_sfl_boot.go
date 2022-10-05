/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package boot

import (
	"github.com/pvillela/go-foa-realworld/internal/config"
	"github.com/pvillela/go-foa-realworld/internal/fl"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
)

///////////////////
// Config logic

var ArticleGetSflCfgAdapter = DefaultSflCfgAdapter

// ArticleGetSflBoot is the function that constructs a stereotype instance of type
// ArticleGetSflT with configuration information and hard-wired stereotype dependencies.
func ArticleGetSflBoot(src config.AppCfgSrc) sfl.ArticleGetSflT {
	return sfl.ArticleGetSflC(
		ArticleGetSflCfgAdapter(src),
		fl.ArticleAndUserGetFl,
	)
}
