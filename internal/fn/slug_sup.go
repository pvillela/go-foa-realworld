package fn

import (
	"github.com/gosimple/slug"
)

func SlugSup(title string) string {
	return slug.Make(title)
}
