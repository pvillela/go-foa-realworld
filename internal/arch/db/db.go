/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package db

// RecCtx is a type that holds platform-specific database record context information, e.g.,
// an optimistic locking token and/or a record ID.
type RecCtx interface {
	Rc() interface{}
}

// Txn defines the abstract interface for transactions.
type Txn interface {

	// Validate returns an error if the transaction is not valid.
	Validate() error

	// End ends the transaction.
	End()
}
