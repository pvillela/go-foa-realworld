/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package bf

import (
	"crypto/ecdsa"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pvillela/go-foa-realworld/experimental/arch/errx"
	"github.com/pvillela/go-foa-realworld/experimental/model"
	"time"
)

type UserGenTokenBfT = func(user model.User) (string, error)

type UserGenTokenEcdsaBfCfgSrc = func() (
	privateKey ecdsa.PrivateKey,
	tokenTimeToLive time.Duration,
)

func UserGenTokenEcdsaBfC(
	cfgSrc UserGenTokenEcdsaBfCfgSrc,
) UserGenTokenBfT {
	privateKey, tokenTimeToLive := cfgSrc()
	return userGenTokenBfC[ecdsa.PrivateKey](privateKey, tokenTimeToLive, jwt.SigningMethodES256)
}

type UserGenTokenHmacBfCfgSrc = func() (
	key []byte,
	tokenTimeToLive time.Duration,
)

func UserGenTokenHmacBfC(
	cfgSrc UserGenTokenHmacBfCfgSrc,
) UserGenTokenBfT {
	key, tokenTimeToLive := cfgSrc()
	return userGenTokenBfC[[]byte](key, tokenTimeToLive, jwt.SigningMethodHS256)
}

func userGenTokenBfC[K any](
	key K,
	tokenTimeToLive time.Duration,
	signingMethod jwt.SigningMethod,
) UserGenTokenBfT {
	return func(user model.User) (string, error) {
		if user.Username == "" {
			return "", errx.NewErrx(nil, "can't generate token for empty user")
		}

		claims := jwt.RegisteredClaims{
			Subject:   user.Username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTimeToLive)),
			Issuer:    "real-world-demo-backend",
		}

		jws, err := jwt.NewWithClaims(signingMethod, &claims).SignedString(key)
		if err != nil {
			return "", errx.ErrxOf(err)
		}

		return jws, nil
	}
}
