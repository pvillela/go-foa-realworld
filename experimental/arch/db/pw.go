/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package db

// Pw wraps a domain entity and a record context together.
// The type parameter E can either be a domain entity type or the pointer type thereof,
// depending on whether the DAF returns / receives by value or by pointer.
// RecCtx holds platform-specific context information about the database record,
// e.g., an optimistic locking token and/or a record Id.
type Pw[E any, R any] struct {
	Entity E `db:""`
	RecCtx R `db:""`
}

// Helper method
func (s Pw[E, R]) Copy(e E) Pw[E, R] {
	s.Entity = e
	return s
}
