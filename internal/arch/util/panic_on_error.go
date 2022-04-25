/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package util

import (
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicLog() {
	if r := recover(); r != nil {
		log.Fatal("panicked:", r.(*errors.Error).ErrorStack())
	}
}
