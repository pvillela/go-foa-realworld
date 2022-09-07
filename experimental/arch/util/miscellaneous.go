/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package util

import (
	"fmt"
	"time"
)

// Ignore is used to remove "unused variable" compilation wrrors without having to remove
// the variable(s).
func Ignore(v ...any) {}

// Trace provides entry and exit tracing for functions, by using it in a defer statement
// as `defer Trace(name)()`
func Trace(name string) func() {
	const timeLayout = "2006-01-02 15:04:05.999999999 Z07:00 MST"
	fmt.Printf("%s entered at %v\n", name, time.Now().Format(timeLayout))
	return func() {
		fmt.Printf("%s returned at %v\n", name, time.Now().Format(timeLayout))
	}
}
