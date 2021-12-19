/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package wgin

import (
	"github.com/pvillela/go-foa-realworld/internal/arch"
	"net/http"

	"github.com/pvillela/go-foa-realworld/internal/arch/web"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func PostHandlerMaker[S any, T any](
	svc func(string, S) (T, error),
	reqCtxExtractor func(*http.Request) (web.RequestContext, error), // TODO: See ExtractToken in jwt_stuff.go
	errorHandler func(arch.Any, web.RequestContext) web.ErrorResult,
) gin.HandlerFunc {

	return func(c *gin.Context) {
		filler := func(pReqCtx *web.RequestContext, pInput *S) error {
			req := c.Request
			reqCtx, err := reqCtxExtractor(req)
			if err != nil {

				log.Info(err)

				c.JSON(401, gin.H{
					"msg": err.Error(),
				})
				return err
			}
			pReqCtx = &reqCtx

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
				return err
			}
			return nil
		}

		setErrorResponse := func(errorContents arch.Any, ctx web.RequestContext) {
			errResult := errorHandler(errorContents, web.RequestContext{})
			c.JSON(errResult.StatusCode, errResult)
		}

		defer func() {
			if r := recover(); r != nil {
				setErrorResponse(r, web.RequestContext{})
			}
		}()

		pseudoHdlr := web.PostPseudoHandler(svc)

		pResp, err := pseudoHdlr(filler)

		if err != nil {
			switch err.(type) {
			case web.FillerError:
				return
			default:
				setErrorResponse(err, web.RequestContext{})
				return
			}
		}

		c.JSON(http.StatusOK, pResp)
	}
}

func SimpleMapGetHandlerMaker(svc func(map[string]string) (arch.Any, error)) gin.HandlerFunc {

	return func(c *gin.Context) {
		params := c.Request.URL.Query()
		m := make(map[string]string, len(params))
		for k, vs := range params {
			m[k] = vs[0]
		}

		pResp, err := svc(m)

		if err != nil {
			c.JSON(400, gin.H{
				"msg":  err.Error(),
				"code": 888,
			})
			return
		}

		c.JSON(http.StatusOK, pResp)
	}
}
