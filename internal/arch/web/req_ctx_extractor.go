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

// DefaultReqCtxExtractor returns a RequestContext based on token if token is not nil,
// returns a zero RequestToken otherwise. It never returns an error.
func DefaultReqCtxExtractor(_ *http.Request, token *jwt.Token) (RequestContext, error) {
	if token == nil {
		return RequestContext{}, nil
	}

	reqCtx := RequestContext{
		Username: token.Claims.(jwt.MapClaims)["sub"].(string),
		Token:    token,
	}
	return reqCtx, nil
}
