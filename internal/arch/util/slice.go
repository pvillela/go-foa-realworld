package util

import (
	"fmt"

	"github.com/pvillela/go-foa-realworld/internal/arch"
)

func SliceWindow(slice []arch.Any, limit, offset int) []arch.Any {
	if slice == nil {
		return []arch.Any{}
	}
	if limit < 0 {
		msg := fmt.Sprintf("SliceWindow limit is %v but should be >= 0")
		panic(msg)
	}
	if offset < 0 {
		msg := fmt.Sprintf("SliceWindow offset is %v but should be >= 0")
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

func SliceContains(slice []arch.Any, value arch.Any) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
