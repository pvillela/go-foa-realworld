/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package util

import (
	"fmt"
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
)

func PanicOnError(err error) {
	if err != nil {
		xErr := errors.New(err)
		panic(xErr)
	}
}

func PanicLog() {
	if r := recover(); r != nil {
		var str string
		switch r.(type) {
		case *errors.Error:
			str = r.(*errors.Error).ErrorStack()
		case error:
			str = r.(error).Error()
		default:
			str = fmt.Sprintf("%v", r)
		}
		log.Fatal("panicked:", str)
	}
}
