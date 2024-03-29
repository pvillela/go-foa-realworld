package boot

import (
	"github.com/pvillela/go-foa-realworld/internal/config"
	"github.com/pvillela/go-foa-realworld/internal/daf"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
)

///////////////////
// Config logic

var ArticleCreateSflCfgAdapter = DefaultSflCfgAdapter

// ArticleCreateSflBoot is the function that constructs a stereotype instance of type
// ArticleCreateSflT with configuration information and hard-wired stereotype dependencies.
func ArticleCreateSflBoot(appCfgSrc config.AppCfgSrc) sfl.ArticleCreateSflT {
	return sfl.ArticleCreateSflC(
		ArticleCreateSflCfgAdapter(appCfgSrc),
		daf.UserGetByNameDaf,
		daf.ArticleCreateDaf,
		daf.TagsAddNewDaf,
		daf.TagsAddToArticleDaf,
	)
}
