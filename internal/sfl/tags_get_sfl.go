/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// TagsGetSflS is the stereotype instance for the service flow that
// retrieves all tags.
type TagsGetSflS struct {
	TagGetAllDaf fs.TagGetAllDafT
}

// TagsGetSflT is the function type instantiated by TagsGetSflS.
type TagsGetSflT = func() (rpc.TagsOut, error)

func (s TagsGetSflS) Make() TagsGetSflT {
	return func() (rpc.TagsOut, error) {
		tags, err := s.TagGetAllDaf()
		if err != nil {
			return rpc.TagsOut{}, err
		}
		tagsOut := rpc.TagsOut_FromModel(tags)
		return tagsOut, err
	}
}
