/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package rpc

type ArticlesListIn struct {
	Tag       string
	Author    string
	Favorited string
	Limit     int
	Offset    int
}
