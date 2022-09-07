/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package errx

// This function is a placeholder for code to be implemented.
// It panics on execution, with an errx.Errx as the panic argument.
// It takes an optional string argument that is used as the message for the aforementioned error object.
// If a message is not provided, the default "missing implementation" is used.
func TODO[T any](args ...string) T {
	var msg string
	if len(args) == 0 {
		msg = "missing implementation"
	} else {
		msg = args[0]
	}
	err := NewErrx(nil, msg)
	panic(err)
}
