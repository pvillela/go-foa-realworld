/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package bf

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"time"
)

type UserGenTokenBfT = func(user model.User) (string, error)

func UserGenTokenBfC(
	tokenTimeToLive time.Duration,
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

		sig, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims).
			SignedString("") // TODO: use appropriate parameter
		if err != nil {
			return "", errx.ErrxOf(err)
		}
		return sig, nil
	}
}
