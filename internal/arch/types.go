/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package arch

type Any = interface{}

type Unit = struct{}

var Void Unit

type Tuple2[T1 any, T2 any] struct {
	_1 T1
	_2 T2
}
