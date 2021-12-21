/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package wgin

import (
	"net/http"

	"github.com/pvillela/go-foa-realworld/internal/arch"

	"github.com/pvillela/go-foa-realworld/internal/arch/web"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type SflT[S any, T any] func(username string, in S) (T, error)

type SflToHandlerT[S any, T any] func(SflT[S, T]) gin.HandlerFunc

type SflToMappedHandlerT[S any, T any] func(
	queryMapper func(map[string]string, *S) error,
	uriMapper func(map[string]string, *S) error,
	svc SflT[S, T],
) gin.HandlerFunc

func setErrorAndAbort(c *gin.Context, err error, httpStatus int) {
	log.Info(err)
	c.JSON(http.StatusBadRequest, gin.H{
		"msg": err.Error(),
	})
	c.AbortWithStatus(httpStatus)
}

// This 3rd order function produces a 2nd order function that takes service flows and returns
// Gin handler functions.
// Contains common logic to bind the HTTP request to service flow input, call service flow, and
// produce HTTP responses.
// - jsonBind: determines whether a JSON payload is expected or not. If true, the request payload
//   will be bound to the service flow input object.
// - queryBind: determines whether query parameters are expected or not. If true, the query params
//   will be bound to the service flow input object.
// - uriBind: determines whether URI parameters are expected or not. If true, the URI params
//   will be bound to the service flow input object.
// - queryMapper: merges a map extracted from query parameters with a service flow input object, to
//   produce an augmented input object.
// - uriMapper: merges a map extracted from uri parameters with a service flow input object, to
//   produce an augmented input object.
func makeSflToHandler[S any, T any](
	jsonBind bool,
	queryBind bool,
	uriBind bool,
	authenticator func(*http.Request) error,
	reqCtxExtractor func(*http.Request) (web.RequestContext, error), // TODO: See ExtractToken in jwt_stuff.go
	errorHandler func(arch.Any, web.RequestContext) web.ErrorResult,
) SflToMappedHandlerT[S, T] {

	return func(
		queryMapper func(map[string]string, *S) error,
		uriMapper func(map[string]string, *S) error,
		svc SflT[S, T],
	) gin.HandlerFunc {

		return func(c *gin.Context) {
			req := c.Request

			err := authenticator(req)
			if err != nil {
				setErrorAndAbort(c, err, 401)
				return
			}

			reqCtx, err := reqCtxExtractor(req)
			if err != nil {
				setErrorAndAbort(c, err, 403) // not sure this is the right code here
				return
			}

			setErrorResponseAndAbort := func(errorContents arch.Any) {
				errResult := errorHandler(errorContents, reqCtx)
				c.JSON(errResult.StatusCode, errResult)
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

			output, err := svc(reqCtx.Username, input)

			if err != nil {
				setErrorResponseAndAbort(err)
				return
			}

			c.JSON(http.StatusOK, &output)
		}
	}
}

func MakeStdNoBodySflToHandler[S any, T any](
	authenticator func(*http.Request) error,
	reqCtxExtractor func(*http.Request) (web.RequestContext, error),
	errorHandler func(arch.Any, web.RequestContext) web.ErrorResult,
) SflToHandlerT[S, T] {
	return func(svc SflT[S, T]) gin.HandlerFunc {
		return makeSflToHandler[S, T](false, true, true, authenticator, reqCtxExtractor, errorHandler)(
			nil, nil, svc,
		)
	}
}

func MakeStdFullBodySflToHandler[S any, T any](
	authenticator func(*http.Request) error,
	reqCtxExtractor func(*http.Request) (web.RequestContext, error),
	errorHandler func(arch.Any, web.RequestContext) web.ErrorResult,
) SflToHandlerT[S, T] {
	return func(svc SflT[S, T]) gin.HandlerFunc {
		return makeSflToHandler[S, T](true, true, true, authenticator, reqCtxExtractor, errorHandler)(
			nil, nil, svc,
		)
	}
}