package daf

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"sync"
)

type TagDafs struct {
	Store *sync.Map
}

func (s TagDafs) MakeGetAll() fs.TagGetAllDafT {
	return func() ([]string, error) {
		var ret []string

		s.Store.Range(func(key, value interface{}) bool {
			tag, ok := key.(string)
			if !ok {
				return true
			}
			ret = append(ret, tag)
			return true
		})

		return ret, nil
	}
}

func (s TagDafs) MakeAdd() fs.TagAddDafT {
	return func(newTags []string) error {
		for _, tag := range newTags {
			s.Store.Store(tag, true)
		}

		return nil
	}
}
