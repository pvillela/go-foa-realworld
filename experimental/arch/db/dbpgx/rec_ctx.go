/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package dbpgx

import "time"

// RecCtx is a type that holds platform-specific database record context information,
// e.g., an optimistic locking token and/or a record Id.  DAFs may accept this type as
// a parameter or return this type, together with domain entity types.
// This type is parameterized to provide type safety, i.e., to prevent passing a RecCtx[U]
// on a call that involves entity type V.
type RecCtx[E any] struct {
	//Id        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
