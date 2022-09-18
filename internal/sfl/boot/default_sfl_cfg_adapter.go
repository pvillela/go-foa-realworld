package boot

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/config"
	"github.com/pvillela/go-foa-realworld/internal/sfl"
)

func defaultSflCfgAdapter(appCfgSrc config.AppCfgInfo) sfl.DefaultSflCfgInfo {
	return util.Todo[sfl.DefaultSflCfgInfo]()
}

var DefaultSflCfgAdapter = util.LiftToNullary(defaultSflCfgAdapter)
