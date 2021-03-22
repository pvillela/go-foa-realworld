/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import "github.com/gin-gonic/gin"

func SetRoutes(r gin.IRouter) {
	r.GET("/tripsvcflow/", tripSvcflowGetH)
	r.POST("/tripsvcflow", tripSvcflowPostH)
}
