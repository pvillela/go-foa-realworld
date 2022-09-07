/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package boot

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/config"
	"github.com/pvillela/go-foa-realworld/internal/fl"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
)

///////////////////
// Config logic

var ArticleGetSflCfgAdapter = func(appCfg config.AppCfgSrc) sfl.DefaultSflCfgSrc {
	return util.Todo[sfl.DefaultSflCfgSrc]()
}

// ArticleGetSflC is the function that constructs a stereotype instance of type
// ArticleGetSflT with hard-wired stereotype dependencies.
func ArticleGetSflC(src config.AppCfgSrc) sfl.ArticleGetSflT {
	return sfl.ArticleGetSflC0(
		ArticleGetSflCfgAdapter(src),
		fl.ArticleAndUserGetFl,
	)
}
