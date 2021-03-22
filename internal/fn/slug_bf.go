package fn

import (
	"github.com/gosimple/slug"
)

func SlugBf(title string) string {
	return slug.Make(title)
}
