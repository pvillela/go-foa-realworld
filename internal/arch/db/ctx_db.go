/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package db

import "context"

type CtxDb interface {
	SetConn(ctx context.Context) (context.Context, error)
	ReleaseConn(ctx context.Context) (context.Context, error)
	DeferredReleaseConn(ctx context.Context)
	BeginTx(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) (context.Context, error)
	Rollback(ctx context.Context) (context.Context, error)
	DeferredRollback(ctx context.Context)
}
