/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package web

func PostPseudoHandler[S any, T any](svc func(S) (T, error)) func(Filler[S]) (T, error) {
	return func(filler Filler[S]) (T, error) {
		var req S
		var resp T
		pReq := &req
		err := filler(pReq)

		if err != nil {
			return resp, FillerError{err}
		}

		resp, err = svc(req)

		return resp, err
	}
}
