package rpc

import "github.com/pvillela/go-foa-realworld/internal/model"

type Profile struct {
	Username  string
	Bio       *string
	Image     *string
	Following bool
}

type ProfileOut struct {
	Profile Profile
}

func (s Profile) FromModel(user model.User, follows bool) Profile {
	if user.Bio != nil {
		s.Bio = user.Bio
	} else {
		empty := ""
		s.Bio = &empty
	}

	if user.ImageLink != "" {
		s.Image = &user.ImageLink
	}

	s.Username = user.Name
	s.Following = follows

	return s
}

func (s ProfileOut) FromModel(user model.User, follows bool) ProfileOut {
	s.Profile = Profile{}.FromModel(user, follows)
	return s
}
