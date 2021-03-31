package rpc

import "github.com/pvillela/go-foa-realworld/internal/model"

type UserOut struct {
	User *userOut0
}

type userOut0 struct {
	Email    string
	Token    string
	Username string
	Bio      string
	Image    *string
}

func (self UserOut) FromModel(user *model.User, token string) UserOut {
	userOut0 := userOut0{
		Email:    user.Email,
		Token:    token,
		Username: user.Name,
		Bio:      *user.Bio,
		Image:    user.ImageLink,
	}
	return UserOut{&userOut0}
}
