package rpc

type TagsOut struct {
	Tags []string
}

func (self TagsOut) FromModel(tags []string) TagsOut {
	return TagsOut{tags}
}
