/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package daftest

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
)

var users = []model.User{
	{
		Username:     "pvillela",
		Email:        "foo@bar.com",
		PasswordHash: "dakfljads0fj",
		PasswordSalt: "2af8d0b50a",
		Bio:          util.PointerFromValue("I am me."),
		ImageLink:    nil,
	},
	{
		Username:     "joebloe",
		Email:        "joe@bloe.com",
		PasswordHash: "9zdakfljads0",
		PasswordSalt: "3ba9e9c611",
		Bio:          util.PointerFromValue("Famous person."),
		ImageLink:    util.PointerFromValue("https://myimage.com"),
	},
	{
		Username:     "johndoe",
		Email:        "johndoe@foo.com",
		PasswordHash: "09fs8asfoasi",
		PasswordSalt: "0000000000",
		Bio:          util.PointerFromValue("Average guy."),
		ImageLink:    util.PointerFromValue("https://johndooeimage.com"),
	},
}

var recCtxUsers = make([]daf.RecCtxUser, len(users))

var articles = []model.Article{
	{
		Title:       "An interesting subject",
		Slug:        "anintsubj",
		Description: "Story about an interesting subject.",
		Body:        util.PointerFromValue("I met this interesting subject a long time ago."),
	},
	{
		Title:       "A dull story",
		Slug:        "adullsubj",
		Description: "Narrative about something dull.",
		Body:        util.PointerFromValue("This is so dull, bla, bla, bla."),
	},
}
