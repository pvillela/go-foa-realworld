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

func MakeDefaultReqCtxExtractor(secretKey []byte) func(*http.Request) (RequestContext, error) {
	return func(req *http.Request) (RequestContext, error) {
		var reqCtx RequestContext

		token, err := VerifiedJwtToken(req, secretKey)
		if err != nil {
			return reqCtx, err
		}

		reqCtx.Username = token.Claims.(jwt.StandardClaims).Subject
		return reqCtx, nil
	}
}
