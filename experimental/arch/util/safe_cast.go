/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package util

// Safely casts x to type T, returning the zero value of T if x is not of type T.
func SafeCast[T any](x any) T {
	var t T
	t, _ = x.(T)
	return t
}
