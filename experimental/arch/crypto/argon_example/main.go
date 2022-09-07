/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/experimental/arch/crypto"
	"github.com/pvillela/go-foa-realworld/experimental/arch/util"
)

func main() {
	defer util.Duration(util.Track("argon2"))
	password := "a fairly lengthy password"
	hash := crypto.ArgonPasswordHash(password)
	fmt.Println(hash)

	hash = crypto.ArgonPasswordHash("password_pvillela")
	fmt.Println(hash)

	hash = crypto.ArgonPasswordHash("password_pvillela")
	fmt.Println(hash)
}
