package rpc

import "github.com/pvillela/go-foa-realworld/internal/model"

type Profile struct {
	Username  string
	Bio       string
	Image     string // image link, nullable
	Following bool
}

func ProfileFromModel(user model.User, following bool) Profile {
	var bio, image string
	if user.Bio != nil {
		bio = *user.Bio
	}
	if user.ImageLink != nil {
		image = *user.ImageLink
	}

	return Profile{
		Username:  user.Name,
		Bio:       bio,
		Image:     image,
		Following: following,
	}
}
