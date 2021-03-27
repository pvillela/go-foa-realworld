package sfl

import (
	"github.com/pvillela/go-foa-realworld/internal/fs"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
)

// TagsGetSfl is the stereotype instance for the service flow that
// retrieves all tags.
type TagsGetSfl struct {
	TagGetAllDaf fs.TagGetAllDafT
}

// TagsGetSflT is the function type instantiated by TagsGetSfl.
type TagsGetSflT = func() (*rpc.TagsOut, error)

func (s TagsGetSfl) Make() TagsGetSflT {
	return func() (*rpc.TagsOut, error) {
		tags, err := s.TagGetAllDaf()
		if err != nil {
			return nil, err
		}
		tagsOut := rpc.TagsOut{}.FromModel(tags)
		return &tagsOut, err
	}
}
