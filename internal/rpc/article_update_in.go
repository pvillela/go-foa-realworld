package rpc

type ArticleUpdateIn struct {
	Article articleUpdateIn0
}

type articleUpdateIn0 struct {
	Title       *string // optional
	Description *string // optional
	Body        *string // optional
}
