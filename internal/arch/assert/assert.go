/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package assert

import (
	"fmt"
	"reflect"
)

func Equal(want interface{}, got interface{}) {
	if !reflect.DeepEqual(got, want) {
		template := `--- Failed equality assertion ...
want:
%#v
got:
%#v
`
		msg := fmt.Sprintf(template, want, got)
		panic(msg)
	}
}

func NotEqual(arg1 interface{}, arg2 interface{}) {
	if reflect.DeepEqual(arg2, arg1) {
		template := `--- Failed inequality assertion ...
arg1:
%#v
arg2:
%#v
`
		msg := fmt.Sprintf(template, arg1, arg2)
		panic(msg)
	}
}

func True(condition bool, msg string) {
	if !condition {
		panic("Assertion failure: " + msg)
	}
}
