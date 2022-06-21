/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import "github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"

// DefaultSflCfgPvdr is the the type of functions that provide
// the required config data for service flow types.
type DefaultSflCfgPvdr = func() (db dbpgx.Db)
