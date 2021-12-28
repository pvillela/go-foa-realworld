/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package web

import (
	"net/http"

	jwt "github.com/golang-jwt/jwt/v4"
)

func DefaultReqCtxExtractor(req *http.Request, token *jwt.Token) (RequestContext, error) {
	reqCtx := RequestContext{
		Username: token.Claims.(jwt.MapClaims)["sub"].(string),
		Token:    token,
	}
	return reqCtx, nil
}
