/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package model

import "time"

type Comment struct {
	Id        uint
	ArticleId uint
	AuthorId  uint
	Body      *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Comment_Create(
	articleId uint,
	authorId uint,
	body *string,
) Comment {
	return Comment{
		Id:        0,
		ArticleId: articleId,
		AuthorId:  authorId,
		Body:      body,
		//CreatedAt: now,
		//UpdatedAt: now,
	}
}
