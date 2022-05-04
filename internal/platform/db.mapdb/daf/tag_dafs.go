/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daf

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/mapdb"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

// TagGetAllDafC is the function that constructs a stereotype instance of type
// bf.TagGetAllDafT.
func TagGetAllDafC(
	tagDb mapdb.MapDb,
) daf.TagGetAllDafT {
	return func() ([]string, error) {
		var ret []string

		tagDb.Store.Range(func(key, value interface{}) bool {
			tag, ok := key.(string)
			if !ok {
				return true
			}
			ret = append(ret, tag)
			return true
		})

		return ret, nil
	}
}

// TagAddDafC is the function that constructs a stereotype instance of type
// bf.TagAddDafT.
func TagAddDafC(
	tagDb mapdb.MapDb,
) daf.TagAddDafT {
	return func(newTags []string) error {
		for _, tag := range newTags {
			tagDb.Store.Store(tag, true)
		}
		return nil
	}
}
