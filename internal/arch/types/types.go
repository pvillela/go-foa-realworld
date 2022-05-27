/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package types

// Unit is the standard functional programming Unit type.
type Unit struct{}

// UnitV is the single instance of Unit.
var UnitV Unit

// CheckType fails to compile if the type doesn't check.
func CheckType[S any](f S) struct{} { return struct{}{} }
