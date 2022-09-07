/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

import "github.com/pvillela/go-foa-realworld/experimental/model"

type TagsOut struct {
	Tags []string
}

func TagsOut_FromModel(tags []model.Tag) TagsOut {
	names := make([]string, len(tags))
	for i := range tags {
		names[i] = tags[i].Name
	}
	return TagsOut{names}
}
