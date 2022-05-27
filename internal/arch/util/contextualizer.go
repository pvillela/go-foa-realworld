/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package util

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
)

type Contextualizer[CC, ERC, IRC, T any] func(
	configCtx CC,
	externalRuntimeCtx ERC,
	block func(externalRuntimeCtx ERC, internalRuntimeCtx IRC) (T, error),
) (T, error)

func LiftContextualizer1[CC, ERC, IRC, S1, T any](
	contextualizer Contextualizer[CC, ERC, IRC, T],
	configCtx CC,
	f func(externalRuntimeCtx ERC, internalRuntimeCtx IRC, s1 S1) (T, error),
) func(externalRuntimeCtx ERC, s1 S1) (T, error) {
	return func(externalRuntimeCtx ERC, s1 S1) (T, error) {
		block := func(externalRuntimeCtx ERC, internalRuntimeCtx IRC) (T, error) {
			return f(externalRuntimeCtx, internalRuntimeCtx, s1)
		}
		return contextualizer(configCtx, externalRuntimeCtx, block)
	}
}

func LiftContextualizer2[CC, ERC, IRC, S1, S2, T any](
	contextualizer Contextualizer[CC, ERC, IRC, T],
	configCtx CC,
	f func(externalRuntimeCtx ERC, internalRuntimeCtx IRC, s1 S1, s2 S2) (T, error),
) func(externalRuntimeCtx ERC, s1 S1, s2 S2) (T, error) {
	return func(externalRuntimeCtx ERC, s1 S1, s2 S2) (T, error) {
		block := func(externalRuntimeCtx ERC, internalRuntimeCtx IRC) (T, error) {
			return f(externalRuntimeCtx, internalRuntimeCtx, s1, s2)
		}
		return contextualizer(configCtx, externalRuntimeCtx, block)
	}
}

func LiftContextualizer1V[CC, ERC, IRC, S1 any](
	contextualizer Contextualizer[CC, ERC, IRC, types.Unit],
	configCtx CC,
	f func(externalRuntimeCtx ERC, internalRuntimeCtx IRC, s1 S1),
) func(externalRuntimeCtx ERC, s1 S1) {
	return func(externalRuntimeCtx ERC, s1 S1) {
		block := func(externalRuntimeCtx ERC, internalRuntimeCtx IRC) (types.Unit, error) {
			f(externalRuntimeCtx, internalRuntimeCtx, s1)
			return types.UnitV, nil
		}
		_, err := contextualizer(configCtx, externalRuntimeCtx, block)
		if err != nil {
			panic(err)
		}
	}
}
