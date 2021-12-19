/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package web

func PostPseudoHandler[S any, T any](
	svc func(string, S) (T, error),
) func(Filler[S]) (T, error) {
	return func(filler Filler[S]) (T, error) {
		var input S
		var output T
		var reqCtx RequestContext
		pReq := &input
		pReqCtx := &reqCtx
		err := filler(pReqCtx, pReq)

		if err != nil {
			return output, FillerError{err}
		}

		output, err = svc(reqCtx.Username, input)

		return output, err
	}
}
