package boot

import (
	"github.com/pvillela/go-foa-realworld/experimental/arch/util"
	"github.com/pvillela/go-foa-realworld/experimental/config"
	"github.com/pvillela/go-foa-realworld/experimental/daf"
	"github.com/pvillela/go-foa-realworld/experimental/sfl"
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
