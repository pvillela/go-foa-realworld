/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

type UserUpdateIn struct {
	User struct {
		Email string
		Bio   *string
		Image string
	}
}
