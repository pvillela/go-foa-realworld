/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/crypto"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
)

func main() {
	defer util.Duration(util.Track("argon2"))
	password := "a fairly lengthy password"
	hash := crypto.BcryptPasswordHash(password)
	fmt.Println(hash)
}
