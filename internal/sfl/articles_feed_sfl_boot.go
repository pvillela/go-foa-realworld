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
)

///////////////////
// Config logic

var ArticlesFeedSflCfgAdapter = func(appCfgSrc config.AppCfgSrc) DefaultSflCfgSrc {
	return util.Todo[DefaultSflCfgSrc]()
}

// ArticlesFeedSflBoot is the function that constructs a stereotype instance of type
// ArticlesFeedSflT with configuration information and hard-wired stereotype dependencies.
func ArticlesFeedSflBoot(src config.AppCfgSrc) ArticlesFeedSflT {
	return ArticlesFeedSflC(
		ArticlesFeedSflCfgAdapter(src),
		daf.UserGetByNameDaf,
		daf.ArticlesFeedDaf,
	)
}
