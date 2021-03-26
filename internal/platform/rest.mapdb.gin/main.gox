/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Create gin router
	router := gin.Default()

	SetRoutes(router)

	// Launch Gin and
	// handle potential error
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
