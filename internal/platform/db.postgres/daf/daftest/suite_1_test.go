/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"testing"
)

func TestSuite1(t *testing.T) {
	dafTester(t, []dbpgx.TestPair{
		{Name: "ArticleDafsSubt", Func: articleDafsSubt},
	})
}
