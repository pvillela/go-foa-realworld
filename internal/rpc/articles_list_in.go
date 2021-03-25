package rpc

type ArticlesListIn struct {
	Tag       string
	Author    string
	Favorited string
	Limit     int
	Offset    int
}
