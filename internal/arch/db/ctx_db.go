/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package db

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
)

type CtxDb interface {
	SetConn(ctx context.Context) (context.Context, error)
	ReleaseConn(ctx context.Context) (context.Context, error)
	DeferredReleaseConn(ctx context.Context)
	BeginTx(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) (context.Context, error)
	Rollback(ctx context.Context) (context.Context, error)
	DeferredRollback(ctx context.Context)
}

func CtxDb_WithTransaction[T any](
	ctxDb CtxDb,
	ctx context.Context,
	block func(ctx context.Context) (T, error),
) (T, error) {
	var zero T
	ctx, err := ctxDb.BeginTx(ctx)
	if err != nil {
		return zero, err
	}

	defer ctxDb.DeferredRollback(ctx)

	t, err := block(ctx)
	if err != nil {
		return zero, err
	}

	_, err = ctxDb.Commit(ctx)
	return t, errx.ErrxOf(err)
}
