/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package olddaf

import (
	"github.com/pvillela/go-foa-realworld/internal/bf"
)

// TagGetAllDafC is the function that constructs a stereotype instance of type
// bf.TagGetAllDafT.
func TagGetAllDafC(
	tagDb mapdb.MapDb,
) bf.TagGetAllDafT {
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
) bf.TagAddDafT {
	return func(newTags []string) error {
		for _, tag := range newTags {
			tagDb.Store.Store(tag, true)
		}
		return nil
	}
}
