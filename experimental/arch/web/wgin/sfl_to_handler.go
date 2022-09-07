/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package wgin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pvillela/go-foa-realworld/experimental/arch/web"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type SflT[S any, T any] func(ctx context.Context, reqCtx web.RequestContext, in S) (T, error)

type SflToHandlerT[S any, T any] func(SflT[S, T]) gin.HandlerFunc

type SflToMappedHandlerT[S any, T any] func(
	queryMapper func(map[string]string, *S) error,
	uriMapper func(map[string]string, *S) error,
	svc SflT[S, T],
) gin.HandlerFunc

func setErrorAndAbort(c *gin.Context, err error, httpStatus int) {
	log.Info(err)
	c.JSON(httpStatus, gin.H{
		"msg": err.Error(),
	})
	c.Abort()
}

// This 3rd order function produces a 2nd order function that takes a service flow and returns
// a Gin handler function.
// Contains common logic to bind the HTTP request to service flow input, call service flow, and
// produce HTTP responses.
//
// - jsonBind: determines whether a JSON payload is expected or not. If true, the request payload
//   will be bound to the service flow input object.
//
// - queryBind: determines whether query parameters are expected or not. If true, the query params
//   will be bound to the service flow input object.
//
// - uriBind: determines whether URI parameters are expected or not. If true, the URI params
//   will be bound to the service flow input object.
//
// - authenticator: authenticates the call, typically using JWT. Is nil for unauthenticated endpoints.
//
// - reqCtxExtractor: extracts information from the HTTP request to form the RequestContext object
//   used throughout the processing flow. If nil, a zero RequestContext is used.
//
// - errorHandler: maps errors before for the HTTP response.
//
// Below are parameters of the function returned by this function:
//
// - queryMapper: merges a map extracted from query parameters with a service flow input object, to
//   produce an augmented input object. Usually is nil.
//
// - uriMapper: merges a map extracted from uri parameters with a service flow input object, to
//   produce an augmented input object. Usually is nil.
//
// - sfl: the service flow that is transformed into a Gin HandlerFunc.
func makeSflHandler[S any, T any](
	jsonBind bool,
	queryBind bool,
	uriBind bool,
	authenticator web.AuthenticatorT,
	authenticationMandatory bool,
	reqCtxExtractor func(*http.Request, *jwt.Token) (web.RequestContext, error),
	errorHandler func(any, web.RequestContext) web.ErrorResult,
) SflToMappedHandlerT[S, T] {

	return func(
		queryMapper func(map[string]string, *S) error,
		uriMapper func(map[string]string, *S) error,
		svc SflT[S, T],
	) gin.HandlerFunc {
		return func(c *gin.Context) {
			req := c.Request

			var token *jwt.Token
			var err error

			// Authentication cases
			rawToken := web.ExtractToken(req)
			switch {
			case (authenticationMandatory || rawToken != "") && authenticator == nil:
				// Application error
				if authenticationMandatory {
					err = fmt.Errorf("authentication mandatory but authenticator is nil")
				} else {
					err = fmt.Errorf("authentication attempt with nil authenticator")
				}
				setErrorAndAbort(c, err, 500)
				return
			case (authenticationMandatory || rawToken != "") && authenticator != nil:
				// Must authenticate
				var ok bool
				ok, token, err = authenticator(req)
				if err != nil {
					setErrorAndAbort(c, err, 401)
					return
				}
				if !ok {
					err = fmt.Errorf("authentication faied")
					setErrorAndAbort(c, err, 401)
					return
				}
			default:
				// token remains nil
			}

			var reqCtx web.RequestContext
			if reqCtxExtractor != nil {
				reqCtx, err = reqCtxExtractor(req, token)
			}
			if err != nil {
				setErrorAndAbort(c, err, 403) // not sure this is the right code here
				return
			}

			setErrorResponseAndAbort := func(errorContents any) {
				errResult := errorHandler(errorContents, reqCtx)
				c.JSON(errResult.StatusCode, errResult)
				c.Abort()
			}

			defer func() {
				if r := recover(); r != nil {
					setErrorResponseAndAbort(r)
				}
			}()

			var input S
			pInput := &input

			if jsonBind {
				// Bind JSON content of request body to pInput
				err = c.ShouldBindJSON(pInput)
				if err != nil {
					setErrorAndAbort(c, err, http.StatusBadRequest)
					return
				}
			}

			if queryBind {
				// Bind query parameters to pInput
				err = c.ShouldBindQuery(pInput)
				if err != nil {
					setErrorAndAbort(c, err, http.StatusBadRequest)
					return
				}
			}

			if uriBind {
				// Bind query parameters to pInput
				err = c.ShouldBindUri(pInput)
				if err != nil {
					setErrorAndAbort(c, err, http.StatusBadRequest)
					return
				}
			}

			if queryMapper != nil {
				params := c.Request.URL.Query()
				m := make(map[string]string, len(params))
				for k, vs := range params {
					m[k] = vs[0]
				}

				err := queryMapper(m, pInput)
				if err != nil {
					if err != nil {
						setErrorAndAbort(c, err, http.StatusBadRequest)
						return
					}
				}
			}

			if uriMapper != nil {
				params := c.Params
				m := make(map[string]string, len(params))
				for _, param := range params {
					m[param.Key] = param.Value
				}

				err := uriMapper(m, pInput)
				if err != nil {
					if err != nil {
						setErrorAndAbort(c, err, http.StatusBadRequest)
						return
					}
				}
			}

			ctx := context.Background()
			output, err := svc(ctx, reqCtx, input)
			if err != nil {
				setErrorResponseAndAbort(err)
				return
			}

			c.JSON(http.StatusOK, output)
		}
	}
}

// MakeLoginHandler is a 3rd order function that produces a 2nd order function that
// takes a login service flow and returns a Gin handler function.
// Contains common logic to bind the HTTP request to service flow input, call service flow, and
// produce HTTP responses.
//
// - errorHandler: maps errors before for the HTTP response.
//
// Below are parameters of the function returned by this function:
//
// - sfl: the service flow that is transformed into a Gin HandlerFunc.
func MakeLoginHandler[S any](
	errorHandler func(any, web.RequestContext) web.ErrorResult,
) SflToHandlerT[S, string] {

	return func(svc SflT[S, string]) gin.HandlerFunc {

		return func(c *gin.Context) {
			setErrorResponseAndAbort := func(errorContents any) {
				errResult := errorHandler(errorContents, web.RequestContext{})
				c.JSON(errResult.StatusCode, errResult)
			}

			defer func() {
				if r := recover(); r != nil {
					setErrorResponseAndAbort(r)
				}
			}()

			var input S
			pInput := &input

			// Bind JSON content of request body to pInput
			err := c.ShouldBindJSON(pInput)
			if err != nil {
				setErrorAndAbort(c, err, http.StatusBadRequest)
				return
			}

			ctx := context.Background()
			output, err := svc(ctx, web.RequestContext{}, input)
			if err != nil {
				setErrorResponseAndAbort(err)
				return
			}

			c.Header("Authorization", output)
			c.Status(http.StatusOK)
		}
	}
}

func MakeStdNoBodySflHandler[S any, T any](
	authenticator func(*http.Request) (bool, *jwt.Token, error),
	authenticationMandatory bool,
	reqCtxExtractor func(*http.Request, *jwt.Token) (web.RequestContext, error),
	errorHandler func(any, web.RequestContext) web.ErrorResult,
) SflToHandlerT[S, T] {
	return func(svc SflT[S, T]) gin.HandlerFunc {
		return makeSflHandler[S, T](
			false,
			true,
			true,
			authenticator,
			authenticationMandatory,
			reqCtxExtractor,
			errorHandler,
		)(
			nil, nil, svc,
		)
	}
}

func MakeStdFullBodySflHandler[S any, T any](
	authenticator web.AuthenticatorT,
	authenticationMandatory bool,
	reqCtxExtractor func(*http.Request, *jwt.Token) (web.RequestContext, error),
	errorHandler func(any, web.RequestContext) web.ErrorResult,
) SflToHandlerT[S, T] {
	return func(svc SflT[S, T]) gin.HandlerFunc {
		return makeSflHandler[S, T](
			true,
			true,
			true,
			authenticator,
			authenticationMandatory,
			reqCtxExtractor,
			errorHandler,
		)(
			nil, nil, svc,
		)
	}
}
