/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package model

import "time"

type Comment struct {
	ArticleId uint
	Id        uint
	AuthorId  uint
	Author    *User
	Body      *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Comment_Create(
	articleId uint,
	body *string,
	author *User,
) Comment {
	now := time.Now()
	comment := Comment{
		ArticleId: articleId,
		Id:        0,
		CreatedAt: now,
		UpdatedAt: now,
		Body:      body,
		Author:    author,
	}
	return comment
}
