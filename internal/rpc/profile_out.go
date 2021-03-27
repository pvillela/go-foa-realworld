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

func (Profile) FromModel(user *model.User, follows bool) Profile {
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
		Following: follows,
	}

	return profile
}

func (ProfileOut) FromModel(user *model.User, follows bool) ProfileOut {
	profile := Profile{}.FromModel(user, follows)
	return ProfileOut{&profile}
}
