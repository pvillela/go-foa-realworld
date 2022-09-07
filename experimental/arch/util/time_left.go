/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package util

import "time"

// TimeLeft converts a deadline to a timeout by subtracting the current time from the deadline.
// If the difference is negative, it returns zero.
func TimeLeft(deadline time.Time) time.Duration {
	timeout := deadline.Sub(time.Now())
	if timeout < 0 {
		return 0
	}
	return timeout
}
