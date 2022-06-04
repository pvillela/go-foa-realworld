package util

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
)

func SliceWindow(slice []any, limit, offset int) []any {
	if slice == nil {
		return []any{}
	}
	if limit < 0 {
		msg := fmt.Sprintf("SliceWindow limit is %v but should be >= 0", limit)
		panic(msg)
	}
	if offset < 0 {
		msg := fmt.Sprintf("SliceWindow offset is %v but should be >= 0", offset)
		panic(msg)
	}
	sLen := len(slice)
	if offset > sLen {
		offset = sLen
	}
	if offset+limit > sLen {
		limit = sLen - offset
	}
	return slice[offset:limit]
}

func SliceContains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// sliceToSet returns a set implemented as a map[T]V that maps all values in the receiver to v.
func sliceToSet[T comparable, V any](s []T, v V) map[T]V {
	if s == nil {
		return nil
	}
	set := make(map[T]V, len(s)) // optimize for speed vs space
	for _, x := range s {
		set[x] = v
	}
	return set
}

// SliceToBoolSet returns a set containing the values in the receiver.
func SliceToBoolSet[T comparable](s []T) map[T]bool {
	return sliceToSet[T, bool](s, true)
}

// SliceToUnitSet returns a set containing the values in the receiver.
func SliceToUnitSet[T comparable](s []T) map[T]types.Unit {
	return sliceToSet[T, types.Unit](s, types.UnitV)
}

func SliceMapWithIndex[S, T any](xs []S, f func(int, S) T) []T {
	if xs == nil {
		return nil
	}
	ts := make([]T, len(xs))
	for i := range xs {
		ts[i] = f(i, xs[i])
	}
	return ts
}

func SliceMap[S, T any](xs []S, f func(S) T) []T {
	if xs == nil {
		return nil
	}
	ts := make([]T, len(xs))
	for i := range xs {
		ts[i] = f(xs[i])
	}
	return ts
}

func SliceFilter[S any](xs []S, f func(S) bool) []S {
	if xs == nil {
		return nil
	}
	ts := make([]S, 0, len(xs))
	for i := range xs {
		if f(xs[i]) {
			ts = append(ts, xs[i])
		}
	}
	return ts
}

func MapToSlice[K comparable, V any](m map[K]V) []V {
	slice := make([]V, len(m))
	i := 0
	for _, v := range m {
		slice[i] = v
		i++
	}
	return slice
}
