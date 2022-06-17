/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfltest

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
	"testing"
	"time"

	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/stretchr/testify/assert"
)

///////////////////
// Constants

const (
	username1 = "pvillela"
	username2 = "joebloe"
	username3 = "johndoe"
)

var secretKey = []byte("abcdefg")

var tokenTimeToLive = func() time.Duration {
	durationStr := "5m"
	dur, err := time.ParseDuration(durationStr)
	errx.PanicOnError(err)
	return dur
}()

///////////////////
// In-memory data

var usertokenMap = make(map[string]string)

///////////////////
// Tests

func userRegisterSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	errx.PanicOnError(err)

	userGenTokenBf := bf.UserGenTokenHmacBfC(secretKey, tokenTimeToLive)
	userRegisterSfl := sfl.UserRegisterSflC(ctxDb, userGenTokenBf)

	userSources := []rpc.UserRegisterIn0{
		{
			Username: username1,
			Email:    "foo@bar.com",
			Password: "password_" + username1,
		},
		{
			Username: username2,
			Email:    "joe@bloe.com",
			Password: "password_" + username2,
		},
		{
			Username: username3,
			Email:    "johndoe@foo.com",
			Password: "password_" + username3,
		},
	}

	badUserSources := []rpc.UserRegisterIn0{
		{ // Existing username
			Username: username1,
			Email:    "foo@bar.com",
			Password: "password_" + username1,
		},
		{ // Existing email
			Username: "dkdkddkdk",
			Email:    "joe@bloe.com",
			Password: "adsklkfjad7809790",
		},
	}

	{
		msg := "user_register_sfl - valid registration"
		for i, _ := range userSources {
			userSrc := userSources[i]
			in := rpc.UserRegisterIn{userSrc}
			out, err := userRegisterSfl(ctx, web.RequestContext{}, in)
			errx.PanicOnError(err)

			// Save tokens in memory
			usertokenMap[in.User.Username] = out.User.Token

			assert.Equal(t, in.User.Username, out.User.Username, msg+" - input Username must match output Username")
			assert.Equal(t, in.User.Email, out.User.Email, msg+" - input Email must match output Email")
		}
	}

	{
		msg := "user_register_sfl - invalid registration"
		for i, _ := range badUserSources {
			userSrc := userSources[i]
			in := rpc.UserRegisterIn{userSrc}
			_, err := userRegisterSfl(ctx, web.RequestContext{}, in)
			returnedErrxKind := dbpgx.ClassifyError(err)
			expectedErrxKind := dbpgx.DbErrUniqueViolation
			expectedErrMsgPrefix := "DbErrUniqueViolation[user with name"

			assert.Equal(t, returnedErrxKind, expectedErrxKind, msg+" - must fail with appropriate error kind when username or email is not unique")
			assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username or email is not unique")
		}
	}
}

//func userAuthenticateSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
//	ctxDb := dbpgx.CtxPgx{db.Pool}
//	ctx, err := ctxDb.SetPool(ctx)
//	errx.PanicOnError(err)
//
//	userGenTokenBf := bf.UserGenTokenHmacBfC(secretKey, tokenTimeToLive)
//	userRegisterSfl := sfl.UserAuthenticateSflC(ctxDb, userGenTokenBf)
//
//	for
//}
