/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package model

type Following struct {
	Id         uint
	FollowerID uint
	//Follower *User
	FolloweeID uint
	//Followee *User
}
