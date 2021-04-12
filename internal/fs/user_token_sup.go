/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package fs

import (
	"errors"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

const tokenTimeToLive = time.Hour * 2

// UserGenTokenSup generates a JWT token for a user
func UserGenTokenSup(user model.User) (string, error) {
	if user.Name == "" {
		return "", errors.New("can't generate token for empty user")
	}

	return jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, newUserClaims(user.Name, tokenTimeToLive)).
		SignedString(user.PasswordSalt)
}

func UserGetNameFromTokenSup(user model.User, tokenStr string) (string, error) {
	token, err := jwtgo.Parse(
		tokenStr,
		func(token *jwtgo.Token) (interface{}, error) {
			// TODO: this function doesn't check the token signature
			return "", nil
		},
	)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwtgo.StandardClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", errors.New("problem with jwt token")
}

// newUserClaims : constructor of userClaims
func newUserClaims(username string, ttl time.Duration) *jwtgo.StandardClaims {
	return &jwtgo.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Issuer:    "real-world-demo-backend",
	}
}
