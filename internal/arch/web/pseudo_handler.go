/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package web

func PostPseudoHandler(pReq Any, svc func(Any) (Any, error)) func(Filler) (Any, error) {
	return func(filler Filler) (Any, error) {
		err := filler(pReq)

		if err != nil {
			return nil, FillerError{err}
		}

		resp, err := svc(pReq)

		return &resp, err
	}
}
