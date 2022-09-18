package sfltest

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/config"
)

func TestCfgAdapterOf[T any](t T) func(config.AppCfgSrc) func() T {
	return util.ConstOf[config.AppCfgSrc, func() T](util.ThunkOf(t))
}
