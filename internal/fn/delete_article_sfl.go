package fn

// DeleteArticleSflS contains the dependencies required for the construction of a
// DeleteArticleSfl. It represents the deletion of an article.
type DeleteArticleSflS struct {
}

// DeleteArticleSfl is the type of a function that takes a slug as input and
// returns noting.
type DeleteArticleSfl = func(slug string)
