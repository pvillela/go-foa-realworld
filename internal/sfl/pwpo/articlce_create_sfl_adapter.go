package pwpo

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/arch/db/dbpgx"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/config"
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
)

///////////////////
// Config logic

func ArticleCreateSflCfgAdapter(appCfg config.AppCfgInfo) sfl.DefaultSflCfgInfo {
	return util.Todo[sfl.DefaultSflCfgInfo]()
}

var ArticleCreateSflCfgSrcV = MakeCfgSrc[sfl.DefaultSflCfgInfo](ArticleCreateSflCfgAdapter)

// ArticleCreateSfl is the configured stereotype instance of type
// ArticleCreateSflT.
var ArticleCreateSfl sfl.ArticleCreateSflT = sfl.ArticleCreateSflC(
	ArticleCreateSflCfgSrcV,
	daf.UserGetByNameDaf,
	daf.ArticleCreateDaf,
	daf.TagsAddNewDaf,
	daf.TagsAddToArticleDaf,
)

var ArticleCreateSfl sfl.ArticleCreateSflT = dbpgx.SflWithTransaction(ArticleCreateSflCfgSrcV(), func(
		ctx context.Context,
		tx pgx.Tx,
		reqCtx web.RequestContext,
		in rpc.ArticleCreateIn,
	) (rpc.ArticleOut, error) {
		err := in.Validate()
		if err != nil {
			return rpc.ArticleOut{}, err
		}
		username := reqCtx.Username

		user, err := daf.UserGetByNameDaf(ctx, tx, username)
		if err != nil {
			return rpc.ArticleOut{}, err
		}

		article, err := in.ToArticle(user)
		if err != nil {
			return rpc.ArticleOut{}, err
		}

		err = daf.ArticleCreateDaf(ctx, tx, &article)
		if err != nil {
			return rpc.ArticleOut{}, err
		}

		names := article.TagList

		err = daf.TagsAddNewDaf(ctx, tx, names)
		if err != nil {
			return rpc.ArticleOut{}, err
		}

		err = daf.TagsAddToArticleDaf(ctx, tx, names, article)
		if err != nil {
			return rpc.ArticleOut{}, err
		}

		articlePlus := model.ArticlePlus_FromArticle(article, false, model.Profile_FromUser(user, false))
		articleOut := rpc.ArticleOut_FromModel(articlePlus)

		return articleOut, nil
	})
}
