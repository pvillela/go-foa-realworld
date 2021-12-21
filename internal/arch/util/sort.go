package util

import (
	"sort"

	"github.com/pvillela/go-foa-realworld/internal/arch"
)

type sortable struct {
	slice []arch.Any
	less  func(i, j int) bool
}

// Len is part of sort.Interface.
func (s sortable) Len() int {
	return len(s.slice)
}

// Swap is part of sort.Interface.
func (s sortable) Swap(i, j int) {
	s.slice[i], s.slice[j] = s.slice[j], s.slice[i]
}

// Less is part of sort.Interface.
func (s sortable) Less(i, j int) bool {
	return s.less(i, j)
}

func Sort(slice []arch.Any, less func(i, j int) bool) {
	s := sortable{slice, less}
	sort.Sort(s)
}
