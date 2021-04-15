package util

import "fmt"

func SliceWindow(slice []Any, limit, offset int) []Any {
	if slice == nil {
		return []Any{}
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

func SliceContains(slice []Any, value Any) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
