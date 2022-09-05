/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfltest

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
)

func makeDefaultSflCfgSrc(db dbpgx.Db) sfl.DefaultSflCfgSrc {
	return func() dbpgx.Db {
		return db
	}
}
