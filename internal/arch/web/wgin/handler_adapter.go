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

type Sfl[S any, T any] func(string, S) (T, error)

type HandlerOfSfl[S any, T any] func(Sfl[S, T]) gin.HandlerFunc

// This 3rd order function produces a 2nd order function that takes service flows and returns
// Gin handler functions.
// Contains common logic to bind the HTTP request to service flow input, call service flow, and
// produce HTTP responses.
// - jsonBind: determines whether a JSON payload is expected or not. If true, the request payload
//   will be bound to the service flow input object.
// - queryMapper: merges a map extracted from query parameters with a service flow input object, to
//   produce an augmented input object.
// - uriMapper: merges a map extracted from uri parameters with a service flow input object, to
//   produce an augmented input object.
func handlerMaker[S any, T any](
	jsonBind bool,
	queryMapper func(map[string]string, *S) error,
	uriMapper func(map[string]string, *S) error,
	authenticator func(*http.Request) error,
	reqCtxExtractor func(*http.Request) (web.RequestContext, error), // TODO: See ExtractToken in jwt_stuff.go
	errorHandler func(arch.Any, web.RequestContext) web.ErrorResult,
) HandlerOfSfl[S, T] {

	return func(svc Sfl[S, T]) gin.HandlerFunc {

		return func(c *gin.Context) {
			req := c.Request

			reqCtx, err := reqCtxExtractor(req)

			setErrorResponse := func(errorContents arch.Any) {
				errResult := errorHandler(errorContents, reqCtx)
				c.JSON(errResult.StatusCode, errResult)
			}

			defer func() {
				if r := recover(); r != nil {
					setErrorResponse(r)
				}
			}()

			if err != nil {

				log.Info(err)

				c.JSON(401, gin.H{
					"msg": err.Error(),
				})
				return
			}

			var input S
			pInput := &input

			if jsonBind {
				// Bind JSON content of request body to pInput
				err = c.BindJSON(pInput)
				if err != nil {
					// Gin automatically returns an error
					// response when the BindJSON operation
					// fails, we simply have to stop this
					// function from continuing to execute

					log.Info(err)

					c.JSON(400, gin.H{
						"msg": err.Error(),
					})
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
					c.JSON(400, gin.H{
						"msg":  err.Error(),
						"code": 888,
					})
					return
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
					c.JSON(400, gin.H{
						"msg":  err.Error(),
						"code": 888,
					})
					return
				}
			}

			output, err := svc(reqCtx.Username, input)

			if err != nil {
				c.JSON(400, gin.H{
					"msg":  err.Error(),
					"code": 888,
				})
				return
			}

			c.JSON(http.StatusOK, &output)
		}
	}
}

func GetHandlerMaker[S any, T any](
	mapper func(map[string]string) (S, error),
	authenticator func(*http.Request) error,
	reqCtxExtractor func(*http.Request) (web.RequestContext, error), // TODO: See ExtractToken in jwt_stuff.go
	errorHandler func(arch.Any, web.RequestContext) web.ErrorResult,
) HandlerOfSfl[S, T] {
	mapper0 := func(m map[string]string, pS *S) error {
		s, err := mapper(m)
		if err != nil {
			return err
		}
		*pS = s
		return nil
	}
	return handlerMaker[S, T](false, mapper0, authenticator, reqCtxExtractor, errorHandler)
}

func MappedBodyHandlerMaker[S any, T any](
	mapper func(map[string]string, *S) error,
	authenticator func(*http.Request) error,
	reqCtxExtractor func(*http.Request) (web.RequestContext, error), // TODO: See ExtractToken in jwt_stuff.go
	errorHandler func(arch.Any, web.RequestContext) web.ErrorResult,
) HandlerOfSfl[S, T] {
	return handlerMaker[S, T](true, mapper, authenticator, reqCtxExtractor, errorHandler)
}

func BodyHandlerMaker[S any, T any](
	authenticator func(*http.Request) error,
	reqCtxExtractor func(*http.Request) (web.RequestContext, error), // TODO: See ExtractToken in jwt_stuff.go
	errorHandler func(arch.Any, web.RequestContext) web.ErrorResult,
) HandlerOfSfl[S, T] {
	return handlerMaker[S, T](true, nil, authenticator, reqCtxExtractor, errorHandler)
}
