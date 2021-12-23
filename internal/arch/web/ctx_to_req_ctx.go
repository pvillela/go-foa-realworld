package web

import (
	"context"

	"github.com/pvillela/go-foa-realworld/internal/arch"
)

func ContextToRequestContext(ctx context.Context) RequestContext {
	return ctx.Value(arch.Void).(RequestContext)
}
