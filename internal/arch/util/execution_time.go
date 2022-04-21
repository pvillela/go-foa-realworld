/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package util

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func Track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func Duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}
