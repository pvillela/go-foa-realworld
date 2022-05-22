/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package arch

// Unit is the standard functional programming Unit type.
type Unit = struct{}

// Void is the single instance of Unit.
var Void Unit

// CheckType fails to compile if the type doesn't check.
func CheckType[S any](f S) struct{} { return struct{}{} }
