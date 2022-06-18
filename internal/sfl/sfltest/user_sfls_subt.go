/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfltest

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
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
// Shared constants and data

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

var userSources = []rpc.UserRegisterIn0{
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

func userAuthenticateSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	errx.PanicOnError(err)

	userGenTokenBf := bf.UserGenTokenHmacBfC(secretKey, tokenTimeToLive)
	userAuthenticateSfl := sfl.UserAuthenticateSflC(ctxDb, userGenTokenBf)

	{
		msg := "user_authenticate_sfl - valid authentication"
		for i, _ := range userSources {
			userSrc := userSources[i]
			in := rpc.UserAuthenticateIn{User: rpc.UserAuthenticateIn0{
				Email:    userSrc.Email,
				Password: userSrc.Password,
			}}

			out, err := userAuthenticateSfl(ctx, web.RequestContext{}, in)
			errx.PanicOnError(err)

			// Save tokens in memory
			usertokenMap[out.User.Username] = out.User.Token

			assert.Equal(t, in.User.Email, out.User.Email, msg+" - input Email must match output Email")
		}
	}

	{
		msg := "user_authenticate_sfl - invalid authentication"

		email := "foo@bar.com"
		password := "abcdefg"

		in := rpc.UserAuthenticateIn{User: rpc.UserAuthenticateIn0{
			Email:    email,
			Password: password,
		}}

		_, err := userAuthenticateSfl(ctx, web.RequestContext{}, in)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := bf.ErrAuthenticationFailed
		expectedErrMsgPrefix := "user authentication failed with name"

		assert.Equal(t, returnedErrxKind, expectedErrxKind, msg+" - must fail with appropriate error kind when username or email is not unique")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username or email is not unique")
	}
}

func userFollowSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	errx.PanicOnError(err)

	userFollowSfl := sfl.UserFollowSflC(ctxDb)

	reqCtx := web.RequestContext{
		Username: username1,
		Token:    &jwt.Token{},
	}

	{
		msg := "user_follow_sfl - follow a valid user not yet followed"

		followeeUsername := username2

		out, err := userFollowSfl(ctx, reqCtx, followeeUsername)
		errx.PanicOnError(err)

		assert.Equal(t, followeeUsername, out.Profile.Username, msg+" - output profile username must equal followee username")
		assert.True(t, out.Profile.Following, msg+" - output profile Following attribute must be true")
	}

	{
		msg := "user_follow_sfl - follow a valid user already followed"

		followeeUsername := username2

		_, err := userFollowSfl(ctx, reqCtx, followeeUsername)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrUniqueViolation
		expectedErrMsgPrefix := "DbErrUniqueViolation[user with username"

		assert.Equal(t, returnedErrxKind, expectedErrxKind, msg+" - must fail with appropriate error kind when followee was already followed")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when followee was already followed")
	}

	{
		msg := "user_follow_sfl - follow an invalid user"

		followeeUsername := "dkdkdkd"

		_, err := userFollowSfl(ctx, reqCtx, followeeUsername)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[user not found for username"

		assert.Equal(t, returnedErrxKind, expectedErrxKind, msg+" - must fail with appropriate error kind when username is not valid")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username is not valid")
	}
}

func userGetCurrentSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	errx.PanicOnError(err)

	userGetCurrentSfl := sfl.UserGetCurrentSflC(ctxDb)

	{
		msg := "user_get_current_sfl - valid username"

		reqCtx := web.RequestContext{
			Username: username1,
			Token:    &jwt.Token{},
		}

		out, err := userGetCurrentSfl(ctx, reqCtx, types.UnitV)
		errx.PanicOnError(err)

		assert.Equal(t, out.User.Username, reqCtx.Username, msg)
	}

	{
		// This test is artificial. In practice, this can never occur due to authentication.

		msg := "user_get_current_sfl - invalid username"

		reqCtx := web.RequestContext{
			Username: "dkdkdkdkd",
			Token:    &jwt.Token{},
		}

		_, err := userGetCurrentSfl(ctx, reqCtx, types.UnitV)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[user not found for username"

		assert.Equal(t, returnedErrxKind, expectedErrxKind, msg+" - must fail with appropriate error kind when username is not valid")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username is not valid")
	}
}
