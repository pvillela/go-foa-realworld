/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package db

// RecCtx is a type that holds platform-specific database record context information,
// e.g., an optimistic locking token and/or a record Id.  DAFs may accept this type as
// a parameter or return this type, together with domain entity types.
// This type is parameterized to provide type safety, i.e., to prevent passing a RecCtx[U]
// on a call that involves entity type V.
type RecCtx[T any] struct {
	Rc interface{}
}

// Pw wraps a domain entity and RecCtx together.  It can be returned or accepted by a
// DAF as an alternative to using RecCtx and the entity type separately.  This is most
// useful when there are multiple entity objects involved as inputs or outputs of a DAF.
// The type parameter T can either be a domain entity type or the pointer type thereof,
// depending on whether the DAF returns / receives by value or by pointer.
type Pw[T any] struct {
	RecCtx[T]
	Entity T
}

// Helper method
func (s Pw[T]) Copy(t T) Pw[T] {
	s.Entity = t
	return s
}

// Txn defines the abstract interface for transactions.
type Txn interface {

	// Validate returns an error if the transaction is not valid.
	Validate() error

	// End ends the transaction.
	End()
}
