package rpc

type TagsOut struct {
	Tags []string
}

func (TagsOut) FromModel(tags []string) TagsOut {
	return TagsOut{tags}
}
