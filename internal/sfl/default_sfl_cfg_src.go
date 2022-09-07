/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/config"
)

// DefaultSflCfgInfo is the type of config data for service flow types.
type DefaultSflCfgInfo = dbpgx.Db

// DefaultSflCfgSrc is the type of functions that provide
// the required config data for service flow types.
type DefaultSflCfgSrc = config.CfgSrc[DefaultSflCfgInfo]
