/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package web

import (
	"fmt"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v4"
)

func DefaultReqCtxExtractor(req *http.Request, claims jwt.MapClaims) (RequestContext, error) {
	var reqCtx RequestContext
	fmt.Println(claims["sub"])
	reqCtx.Username = claims["sub"].(string)
	return reqCtx, nil
}
