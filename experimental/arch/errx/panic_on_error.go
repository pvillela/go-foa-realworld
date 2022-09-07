/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package errx

import (
	"fmt"
	"runtime/debug"
)

func PanicOnError(err error) {
	if err != nil {
		panic(ErrxOf(err))
	}
}

func PanicLog(logger func(args ...interface{})) {
	if r := recover(); r != nil {
		var str string
		switch r.(type) {
		case Errx:
			str = r.(Errx).StackTrace()
		default:
			var errStr string
			errStr = fmt.Sprintf("%v", r)
			stack := debug.Stack()
			str = errStr + "\n" + string(stack)
		}
		logger("panicked: ", str)
	}
}
