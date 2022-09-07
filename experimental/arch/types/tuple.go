/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package types

// Tuple2 is tuple with 2 elements
type Tuple2[T1, T2 any] struct {
	X1 T1
	X2 T2
}

// Tuple3 is tuple with 3 elements
type Tuple3[T1, T2, T3 any] struct {
	X1 T1
	X2 T2
	X3 T3
}

// Tuple4 is tuple with 4 elements
type Tuple4[T1, T2, T3, T4 any] struct {
	X1 T1
	X2 T2
	X3 T3
	X4 T4
}
