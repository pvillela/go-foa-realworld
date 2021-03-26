package rpc

import "github.com/pvillela/go-foa-realworld/internal/model"

type Profile struct {
	Username  string
	Bio       string
	Image     string
	Following bool
}

type ProfileOut struct {
	Profile *Profile
}

func (Profile) FromModel(user model.User, following bool) Profile {
	var bio, image string
	if user.Bio != nil {
		bio = *user.Bio
	}
	if user.ImageLink != nil {
		image = *user.ImageLink
	}

	profile := Profile{
		Username:  user.Name,
		Bio:       bio,
		Image:     image,
		Following: following,
	}

	return profile
}

func (ProfileOut) FromModel(user model.User, following bool) ProfileOut {
	profile := Profile{}.FromModel(user, following)
	return ProfileOut{&profile}
}
