/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfltest

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/cdb"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
	rpc2 "github.com/pvillela/go-foa-realworld/rpc"
	"testing"
	"time"

	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/stretchr/testify/assert"
)

///////////////////
// Shared constants and data

const (
	username1 = "pvillela"
	username2 = "joebloe"
	username3 = "johndoe"
	username4 = "initial_username4"
)

var secretKey = []byte("abcdefg")

var tokenTimeToLive = func() time.Duration {
	durationStr := "5m"
	dur, err := time.ParseDuration(durationStr)
	errx.PanicOnError(err)
	return dur
}()

var userSources = map[string]rpc2.UserRegisterIn0{
	username1: {
		Username: username1,
		Email:    "foo@bar.com",
		Password: "password_" + username1,
	},
	username2: {
		Username: username2,
		Email:    "joe@bloe.com",
		Password: "password_" + username2,
	},
	username3: {
		Username: username3,
		Email:    "johndoe@foo.com",
		Password: "password_" + username3,
	},
	username4: {
		Username: username4,
		Email:    username4 + "@foo.com",
		Password: "password_" + username4,
	},
}

///////////////////
// Helpers

func makeUserGenTokenHmacBfCfgPvdr(key []byte, tokenTtl time.Duration) bf.UserGenTokenHmacBfCfgPvdr {
	return func() ([]byte, time.Duration) {
		return key, tokenTtl
	}
}

func makeUserSflCfgPvdr(ctxDb cdb.CtxDb) sfl.UserSflCfgPvdr {
	return func() cdb.CtxDb {
		return ctxDb
	}
}

///////////////////
// Tests

func userRegisterSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	assert.NoError(t, err)

	userGenTokenBf := bf.UserGenTokenHmacBfC(makeUserGenTokenHmacBfCfgPvdr(secretKey, tokenTimeToLive))
	userRegisterSfl := sfl.UserRegisterSflC(makeUserSflCfgPvdr(ctxDb), userGenTokenBf)

	{
		msg := "user_register_sfl - valid registration"
		for k, _ := range userSources {
			userSrc := userSources[k]
			in := rpc2.UserRegisterIn{userSrc}
			out, err := userRegisterSfl(ctx, web.RequestContext{}, in)
			assert.NoError(t, err)

			assert.Equal(t, in.User.Username, out.User.Username, msg+" - input Username must match output Username")
			assert.Equal(t, in.User.Email, out.User.Email, msg+" - input Email must match output Email")
		}
	}

	{
		msg := "user_register_sfl - invalid registration"

		badUserSources := []rpc2.UserRegisterIn0{
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
			userSrc := badUserSources[i]
			in := rpc2.UserRegisterIn{userSrc}
			_, err := userRegisterSfl(ctx, web.RequestContext{}, in)
			returnedErrxKind := dbpgx.ClassifyError(err)
			expectedErrxKind := dbpgx.DbErrUniqueViolation
			expectedErrMsgPrefix := "DbErrUniqueViolation[user with name"

			assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when username or email is not unique")
			assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username or email is not unique")
		}
	}
}

func userAuthenticateSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	assert.NoError(t, err)

	userGenTokenBf := bf.UserGenTokenHmacBfC(makeUserGenTokenHmacBfCfgPvdr(secretKey, tokenTimeToLive))
	userAuthenticateSfl := sfl.UserAuthenticateSflC(makeUserSflCfgPvdr(ctxDb), userGenTokenBf)

	{
		msg := "user_authenticate_sfl - valid authentication"
		for k, _ := range userSources {
			userSrc := userSources[k]
			in := rpc2.UserAuthenticateIn{User: rpc2.UserAuthenticateIn0{
				Email:    userSrc.Email,
				Password: userSrc.Password,
			}}

			out, err := userAuthenticateSfl(ctx, web.RequestContext{}, in)
			assert.NoError(t, err)

			assert.Equal(t, in.User.Email, out.User.Email, msg+" - input Email must match output Email")
		}
	}

	{
		msg := "user_authenticate_sfl - invalid authentication"

		email := "foo@bar.com"
		password := "abcdefg"

		in := rpc2.UserAuthenticateIn{User: rpc2.UserAuthenticateIn0{
			Email:    email,
			Password: password,
		}}

		_, err := userAuthenticateSfl(ctx, web.RequestContext{}, in)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := bf.ErrAuthenticationFailed
		expectedErrMsgPrefix := "user authentication failed with name"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when username or email is not unique")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username or email is not unique")
	}
}

func userFollowSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	assert.NoError(t, err)

	userFollowSfl := sfl.UserFollowSflC(makeUserSflCfgPvdr(ctxDb))

	reqCtx := web.RequestContext{
		Username: username1,
		Token:    &jwt.Token{},
	}

	{
		msg := "user_follow_sfl - follow a valid user not yet followed"

		followeeUsername := username2

		out, err := userFollowSfl(ctx, reqCtx, followeeUsername)
		assert.NoError(t, err)

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

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when followee was already followed")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when followee was already followed")
	}

	{
		msg := "user_follow_sfl - follow an invalid user"

		followeeUsername := "dkdkdkd"

		_, err := userFollowSfl(ctx, reqCtx, followeeUsername)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[user not found for username"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when username is not valid")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username is not valid")
	}
}

func userGetCurrentSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	assert.NoError(t, err)

	userGetCurrentSfl := sfl.UserGetCurrentSflC(makeUserSflCfgPvdr(ctxDb))

	{
		msg := "user_get_current_sfl - valid username"

		reqCtx := web.RequestContext{
			Username: username1,
			Token:    &jwt.Token{},
		}

		out, err := userGetCurrentSfl(ctx, reqCtx, types.UnitV)
		assert.NoError(t, err)

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

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when username is not valid")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username is not valid")
	}
}

func userUnfollowSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	assert.NoError(t, err)

	userFollowSfl := sfl.UserUnfollowSflC(makeUserSflCfgPvdr(ctxDb))

	reqCtx := web.RequestContext{
		Username: username1,
		Token:    &jwt.Token{},
	}

	{
		msg := "user_unfollow_sfl - unfollow a valid user currently followed"

		followeeUsername := username2

		out, err := userFollowSfl(ctx, reqCtx, followeeUsername)
		assert.NoError(t, err)

		assert.Equal(t, followeeUsername, out.Profile.Username, msg+" - output profile username must equal followee username")
		assert.False(t, out.Profile.Following, msg+" - output profile Following attribute must be false")
	}

	{
		msg := "user_unfollow_sfl - unfollow a valid user not already followed"

		followeeUsername := username2

		_, err := userFollowSfl(ctx, reqCtx, followeeUsername)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[user with username"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when followee was not already followed")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when followee was not already followed")
	}

	{
		msg := "user_unfollow_sfl - unfollow an invalid user"

		followeeUsername := "dkdkdkd"

		_, err := userFollowSfl(ctx, reqCtx, followeeUsername)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrRecordNotFound
		expectedErrMsgPrefix := "DbErrRecordNotFound[user not found for username"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when username is not valid")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when username is not valid")
	}
}

func userUpdateSflSubt(db dbpgx.Db, ctx context.Context, t *testing.T) {
	ctxDb := dbpgx.CtxPgx{db.Pool}
	ctx, err := ctxDb.SetPool(ctx)
	assert.NoError(t, err)

	userUpdateSfl := sfl.UserUpdateSflC(makeUserSflCfgPvdr(ctxDb))

	reqCtx := web.RequestContext{
		Username: username4,
		Token:    &jwt.Token{},
	}

	{
		msg := "user_update_sfl - valid changes, same username and email"

		updatedUsername := username4
		updatedEmail := userSources[updatedUsername].Email

		in := rpc2.UserUpdateIn{User: model.UserPatch{
			Username:  &updatedUsername,
			Email:     &updatedEmail,
			Password:  util.PointerFromValue("password_" + updatedUsername),
			Bio:       util.PointerFromValue("I am the 4th user."),
			ImageLink: util.PointerFromValue("http://foo.com/" + updatedUsername + ".png"),
		}}

		out, err := userUpdateSfl(ctx, reqCtx, in)
		assert.NoError(t, err)

		expected := rpc2.UserOut{User: rpc2.UserOut0{
			Email:    *in.User.Email,
			Token:    reqCtx.Token.Raw,
			Username: *in.User.Username,
			Bio:      in.User.Bio,
			Image:    *in.User.ImageLink,
		}}

		assert.Equal(t, expected, out, msg+" - output must align with changes")
	}

	{
		msg := "user_update_sfl - duplicate username"

		updatedUsername := username1

		in := rpc2.UserUpdateIn{User: model.UserPatch{
			Username:  &updatedUsername,
			Email:     util.PointerFromValue(updatedUsername + "@foo.com"),
			Password:  util.PointerFromValue("password_" + updatedUsername),
			Bio:       util.PointerFromValue("I am the 4th user."),
			ImageLink: util.PointerFromValue("http://foo.com/" + updatedUsername + ".png"),
		}}

		_, err := userUpdateSfl(ctx, reqCtx, in)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrUniqueViolation
		expectedErrMsgPrefix := "DbErrUniqueViolation[user with name"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when new username already exists")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when new username already exists")
	}

	{
		msg := "user_update_sfl - duplicate email"

		updatedUsername := username4

		in := rpc2.UserUpdateIn{User: model.UserPatch{
			Username:  &updatedUsername,
			Email:     util.PointerFromValue(userSources[username1].Email),
			Password:  util.PointerFromValue("password_" + updatedUsername),
			Bio:       util.PointerFromValue("I am the 4th user."),
			ImageLink: util.PointerFromValue("http://foo.com/" + updatedUsername + ".png"),
		}}

		_, err := userUpdateSfl(ctx, reqCtx, in)
		returnedErrxKind := dbpgx.ClassifyError(err)
		expectedErrxKind := dbpgx.DbErrUniqueViolation
		expectedErrMsgPrefix := "DbErrUniqueViolation[user with name"

		assert.Equal(t, expectedErrxKind, returnedErrxKind, msg+" - must fail with appropriate error kind when new email already exists")
		assert.ErrorContains(t, err, expectedErrMsgPrefix, msg+" - must fail with appropriate error message when new email already exists")
	}

	{
		msg := "user_update_sfl - valid changes, different username and email"

		updatedUsername := "updated_username4"

		in := rpc2.UserUpdateIn{User: model.UserPatch{
			Username:  &updatedUsername,
			Email:     util.PointerFromValue(updatedUsername + "@foo.com"),
			Password:  util.PointerFromValue("password_" + updatedUsername),
			Bio:       util.PointerFromValue("I am the 4th user."),
			ImageLink: util.PointerFromValue("http://foo.com/" + updatedUsername + ".png"),
		}}

		out, err := userUpdateSfl(ctx, reqCtx, in)
		assert.NoError(t, err)

		expected := rpc2.UserOut{User: rpc2.UserOut0{
			Email:    *in.User.Email,
			Token:    reqCtx.Token.Raw,
			Username: *in.User.Username,
			Bio:      in.User.Bio,
			Image:    *in.User.ImageLink,
		}}

		assert.Equal(t, expected, out, msg+" - output must align with changes")
	}
}
