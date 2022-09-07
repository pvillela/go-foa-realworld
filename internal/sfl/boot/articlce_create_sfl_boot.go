package boot

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/config"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
)

///////////////////
// Config logic

var ArticleCreateSflCfgAdapter = func(appCfgSrc config.AppCfgSrc) sfl.DefaultSflCfgSrc {
	return util.Todo[sfl.DefaultSflCfgSrc]()
}

// ArticleCreateSflBoot is the function that constructs a stereotype instance of type
// ArticleCreateSflT with configuration information and hard-wired stereotype dependencies.
func ArticleCreateSflBoot(appCfgSrc config.AppCfgSrc) sfl.ArticleCreateSflT {
	return sfl.ArticleCreateSflC0(
		ArticleCreateSflCfgAdapter(appCfgSrc),
		daf.UserGetByNameExplicitTxDaf,
		daf.ArticleCreateDaf,
		daf.TagsAddNewDaf,
		daf.TagsAddToArticleDaf,
	)
}
