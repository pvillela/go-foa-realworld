/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/config"
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/sfl/boot"
)

///////////////////
// Config logic

// ArticlesListSflBoot is the function that constructs a stereotype instance of type
// ArticlesListSflT with configuration information and hard-wired stereotype dependencies.
func ArticlesListSflBoot(src config.AppCfgSrc) ArticlesListSflT {
	return ArticlesListSflC(
		boot.DefaultSflCfgAdapter(src),
		daf.UserGetByNameDaf,
		daf.ArticlesListDaf,
	)
}
