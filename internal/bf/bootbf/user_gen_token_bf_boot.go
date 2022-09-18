package bootbf

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/config"
)

///////////////////
// Config logic

var userGenTokenBfCfgAdapter = func(appCfgSrc config.AppCfgInfo) bf.UserGenTokenHmacBfCfgInfo {
	return util.Todo[bf.UserGenTokenHmacBfCfgInfo]()
}

var UserGenTokenBfCfgAdapter = util.LiftToNullary(userGenTokenBfCfgAdapter)

func UserGenTokenHmacBfBoot(appCfgSrc config.AppCfgSrc) bf.UserGenTokenBfT {
	src := UserGenTokenBfCfgAdapter(appCfgSrc)
	return bf.UserGenTokenHmacBfC(src)
}
