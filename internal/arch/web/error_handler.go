/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package web

import (
	"github.com/sirupsen/logrus"
)

func DefaultErrorHandler(errorContents any, ctx RequestContext) ErrorResult {
	logrus.Info("Error caught in web handler:", errorContents)
	var errResult ErrorResult
	errResult.DeveloperMessage = "Dummy error handler implementation."
	errResult.StatusCode = 500
	switch errorContents.(type) {
	case error:
		errResult.ErrorString = errorContents.(error).Error()
	}
	return errResult
}
