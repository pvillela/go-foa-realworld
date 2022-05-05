/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package util

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
)

func PanicOnError(err error) {
	if err != nil {
		panic(errx.ErrxOf(err))
	}
}

func PanicLog(logger func(args ...interface{})) {
	if r := recover(); r != nil {
		var str string
		switch r.(type) {
		case errx.Errx:
			str = fmt.Sprintf("%+v", r)
		case *errors.Error:
			str = r.(*errors.Error).ErrorStack()
		case error:
			str = r.(error).Error()
		default:
			str = fmt.Sprintf("%+v", r)
		}
		logger("panicked:", str)
	}
}
