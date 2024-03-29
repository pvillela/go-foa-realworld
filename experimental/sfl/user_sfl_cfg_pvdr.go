/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/experimental/arch/db/cdb"
)

// UserSflCfgSrc is the the type of functions that provide
// the required config data for User service flow types.
type UserSflCfgSrc = func() (ctxDb cdb.CtxDb)
