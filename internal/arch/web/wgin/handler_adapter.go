/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package wgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvillela/gfoa/pkg/web"
	log "github.com/sirupsen/logrus"
)

func PostHanderMaker(
	pInput Any,
	svc func(Any) (Any, error),
	errorHandler func(Any, web.RequestContext) web.ErrorResult,
) gin.HandlerFunc {

	return func(c *gin.Context) {
		filler := func(pInput Any) error {
			// Bind JSON content of request body to
			// struct created above
			err := c.BindJSON(pInput)
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

		setErrorResponse := func(errorContents Any, ctx web.RequestContext) {
			errResult := errorHandler(errorContents, web.RequestContext{})
			c.JSON(errResult.StatusCode, errResult)
		}

		defer func() {
			if r := recover(); r != nil {
				setErrorResponse(r, web.RequestContext{})
			}
		}()

		pseudoHdlr := web.PostPseudoHandler(pInput, svc)

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

func SimpleMapGetHanderMaker(svc func(map[string]string) (Any, error)) gin.HandlerFunc {

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
