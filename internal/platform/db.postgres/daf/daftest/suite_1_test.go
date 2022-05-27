/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"testing"
)

func TestSuite1(t *testing.T) {
	dafTester0(t, []TestPair0{
		{Name: "ArticleDafsSubt", Func: articleDafsSubt0},
	})

	dafTester1(t, []TestPair0{
		{Name: "ArticleDafsSubt", Func: articleDafsSubt1},
	})
}
