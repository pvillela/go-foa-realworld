/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// TagsGetSflT is the type of the stereotype instance for the service flow that
// retrieves all tags.
type TagsGetSflT = func(_ context.Context, _ arch.Unit) (rpc.TagsOut, error)

// TagsGetSflC is the function that constructs a stereotype instance of type
// TagsGetSflT.
func TagsGetSflC(
	tagGetAllDaf daf.TagGetAllDafT,
) TagsGetSflT {
	return func(_ context.Context, _ arch.Unit) (rpc.TagsOut, error) {
		tags, err := tagGetAllDaf()
		if err != nil {
			return rpc.TagsOut{}, err
		}
		tagsOut := rpc.TagsOut_FromModel(tags)
		return tagsOut, err
	}
}
