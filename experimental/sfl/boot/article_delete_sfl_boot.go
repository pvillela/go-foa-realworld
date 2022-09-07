/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package boot

import (
	"github.com/pvillela/go-foa-realworld/experimental/arch/util"
	"github.com/pvillela/go-foa-realworld/experimental/config"
	"github.com/pvillela/go-foa-realworld/experimental/fl"
	"github.com/pvillela/go-foa-realworld/experimental/daf"
	"github.com/pvillela/go-foa-realworld/experimental/sfl"
)

///////////////////
// Config logic

var ArticleDeleteSflCfgAdapter = func(appCfgSrc config.AppCfgSrc) sfl.DefaultSflCfgSrc {
	return util.Todo[sfl.DefaultSflCfgSrc]()
}

// ArticleDeleteSflC is the function that constructs a stereotype instance of type
// ArticleDeleteSflT with hard-wired stereotype dependencies.
func ArticleDeleteSflC(appCfgSrc config.AppCfgSrc) sfl.ArticleDeleteSflT {
	return sfl.ArticleDeleteSflC0(
		ArticleDeleteSflCfgAdapter(appCfgSrc),
		fl.ArticleGetAndCheckOwnerFl,
		daf.ArticleDeleteDaf,
	)
}
