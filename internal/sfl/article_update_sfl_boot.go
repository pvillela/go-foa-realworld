/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/config"
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/fl"
)

///////////////////
// Config logic

var ArticleUpdateSflCfgAdapter = func(appCfgSrc config.AppCfgSrc) DefaultSflCfgSrc {
	return util.Todo[DefaultSflCfgSrc]()
}

// ArticleUpdateSflBoot is the function that constructs a stereotype instance of type
// ArticleUpdateSflT with configuration information and hard-wired stereotype dependencies.
func ArticleUpdateSflBoot(src config.AppCfgSrc) ArticleUpdateSflT {
	return ArticleUpdateSflC0(
		ArticleUpdateSflCfgAdapter(src),
		fl.ArticleGetAndCheckOwnerFl,
		daf.ArticleUpdateDaf,
	)
}
