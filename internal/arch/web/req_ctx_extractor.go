/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package web

import "net/http"

func DefaultReqCtxExtractor(*http.Request) (RequestContext, error) {
	panic("todo")
}
