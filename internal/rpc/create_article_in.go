package rpc

type CreateArticleIn struct {
	Article struct {
		Title       string
		Description string
		Body        string
		TagList     []string // optional
	}
}
