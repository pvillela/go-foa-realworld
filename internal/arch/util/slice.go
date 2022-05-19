package util

import (
	"fmt"
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
