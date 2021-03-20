package rpc

type UpdateArticleIn struct {
	Article struct {
		Title       string // optional
		Description string // optional
		Body        string // optional
	}
}
