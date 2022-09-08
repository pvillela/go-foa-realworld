package bootbf

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/bf"
	"github.com/pvillela/go-foa-realworld/internal/config"
)

///////////////////
// Config logic

var UserGenTokenBfCfgAdapter = func(appCfgSrc config.AppCfgSrc) bf.UserGenTokenHmacBfCfgSrc {
	return util.Todo[bf.UserGenTokenHmacBfCfgSrc]()
}

func UserGenTokenHmacBfBoot(appCfgSrc config.AppCfgSrc) bf.UserGenTokenBfT {
	src := UserGenTokenBfCfgAdapter(appCfgSrc)
	return bf.UserGenTokenHmacBfC(src)
}
