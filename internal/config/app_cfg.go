package config

import "github.com/pvillela/go-foa-realworld/internal/arch/util"

type AppCfgInfo struct {
}

type AppCfgSrc = func() AppCfgInfo

func GetAppConfiguration() AppCfgInfo {
	util.Todo[any]()
	return AppCfgInfo{}
}
