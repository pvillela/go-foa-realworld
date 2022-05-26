package util

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch"
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

func SliceContains(slice []any, value any) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// SliceToSet returns a set containing the values in the receiver.
func SliceToSet[T comparable](s []T) map[T]arch.Unit {
	if s == nil {
		return nil
	}
	set := make(map[T]arch.Unit, len(s)) // optimize for speed vs space
	for _, x := range s {
		set[x] = arch.Void
	}
	return set
}

func SliceMap[S, T any](xs []S, f func(int, S) T) []T {
	if xs == nil {
		return nil
	}
	ts := make([]T, len(xs))
	for i := range xs {
		ts[i] = f(i, xs[i])
	}
	return ts
}
