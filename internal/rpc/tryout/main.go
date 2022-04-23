/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

func main() {
	bio := "I am a Bar."
	user := model.User{
		Username:       "Foo",
		Email:          "foo@bar.com",
		IsTempPassword: false,
		PasswordHash:   "xhxh",
		PasswordSalt:   nil,
		Bio:            &bio,
		ImageLink:      "",
		Followees:      nil,
		//Favorites:      nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Add(10),
	}

	userOut := rpc.UserOut_FromModel(user, "abc")

	bytes, err := json.Marshal(userOut)
	str := string(bytes)
	fmt.Println(err, str)
}
