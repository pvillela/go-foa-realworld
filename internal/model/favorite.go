/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package model

type Favorite struct {
	// UserId is the ID of the user that favors the article.
	UserId    uint
	ArticleId uint
}
