/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package cdb

import (
	"context"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
)

type CtxDb interface {
	BeginTx(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) (context.Context, error)
	Rollback(ctx context.Context) (context.Context, error)
	DeferredRollback(ctx context.Context)
}

func WithTransaction[T any](
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

func SflWithTransaction[R, S, T any](
	ctxDb CtxDb,
	block func(ctx context.Context, reqCtx R, in S) (T, error),
) func(ctx context.Context, reqCtx R, in S) (T, error) {
	return func(ctx context.Context, reqCtx R, in S) (T, error) {
		block1 := func(ctx context.Context) (T, error) {
			return block(ctx, reqCtx, in)
		}
		return WithTransaction(ctxDb, ctx, block1)
	}
}
