/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package web

import "github.com/sirupsen/logrus"

func DefaultErrorHandler(errorContents Any, ctx RequestContext) ErrorResult {
	logrus.Info("Error caught in web handler:", errorContents)
	errResult := ErrorResult{}
	errResult.DeveloperMessage = "Dummy error handler implementation."
	errResult.StatusCode = 500
	return errResult
}
