/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package model

// ArticleCriteria defines criteria for selection of articles
type ArticleCriteria struct {
	// All fields are optional
	Tag         *string `json:"tag,omitempty"`
	Author      *string `json:"author,omitempty"`
	FavoritedBy *string `json:"favorited,omitempty"`
	Limit       *int    `json:"limit,omitempty"`  // default 20
	Offset      *int    `json:"offset,omitempty"` // default 0
}
