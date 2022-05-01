/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package sfl

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/newdaf"

	"github.com/pvillela/go-foa-realworld/internal/arch"
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/fs"
)

// ArticleDeleteSflT is the type of the stereotype instance for the service flow that
// deletes an article.
type ArticleDeleteSflT = func(ctx context.Context, slug string) (arch.Unit, error)

// ArticleDeleteSflC is the function that constructs a stereotype instance of type
// ArticleDeleteSflT.
func ArticleDeleteSflC(
	beginTxn func(context string) db.Txn,
	articleGetAndCheckOwnerFl fs.ArticleGetAndCheckOwnerFlT,
	articleDeleteDaf newdaf.ArticleDeleteDafT,
) ArticleDeleteSflT {
	return func(ctx context.Context, slug string) (arch.Unit, error) {
		username := web.ContextToRequestContext(ctx).Username
		txn := beginTxn("ArticleCreateSflS")
		defer txn.End()

		_, _, err := articleGetAndCheckOwnerFl(username, slug)
		if err != nil {
			return arch.Void, err
		}

		return arch.Void, articleDeleteDaf(slug, txn)
	}
}
