package config

type AppCfgInfo struct {
}

type AppCfgSrc func() AppCfgInfo

func GetAppConfiguration() AppCfgInfo {
	return AppCfgInfo{}
}
