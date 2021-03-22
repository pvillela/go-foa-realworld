package rpc

type ArticleUpdateIn struct {
	Article ArticleUpdateIn0
}

type ArticleUpdateIn0 struct {
	Title       string // optional
	Description string // optional
	Body        string // optional
}
