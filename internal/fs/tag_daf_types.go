/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package fs

// TagGetAllDafT is the type of the stereotype instance for the DAF that
// retrieves all tags.
type TagGetAllDafT = func() ([]string, error)

// TagAddDafT is the type of the stereotype instance for the DAF that
// adds a tag.
type TagAddDafT = func(newTags []string) error
