/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package web

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AuthenticatorT = func(*http.Request) (bool, *jwt.Token, error)

// AuthenticatorC constructs an AuthenticatorT function from the secretKey argument.
// The constructed authenticator only uses the JWT token signature and expiration to
// authenticate. A more elaborate authenticator could use a cache to enable immediate
// token invalidation in case of logout.
func AuthenticatorC(
	secretKey []byte,
) AuthenticatorT {
	return func(req *http.Request) (bool, *jwt.Token, error) {
		token, err := VerifiedJwtToken(req, secretKey)
		if err != nil {
			return false, token, err
		}
		logrus.Debug("authenticator ran\n", "token", token)
		return true, token, nil
	}
}
