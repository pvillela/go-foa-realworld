package sfl

// DeleteCommentSflS contains the dependencies required for the construction of a
// DeleteCommentSfl. It represents the deletion of a comment to an article.
type DeleteCommentSflS struct {
}

// DeleteCommentSfl is the type of a function that takes a slug and a comment id as inputs
// and returns nothing.
type DeleteCommentSfl = func(slug string, commentId string)
