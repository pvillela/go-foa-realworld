package main

import (
	"encoding/json"
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/rpc"
	"time"
)

func main() {
	bio := "I am a Bar."
	user := model.User{
		Name:         "Foo",
		Email:        "foo@bar.com",
		TempPassword: "temp",
		PasswordHash: "xhxh",
		PasswordSalt: nil,
		Bio:          &bio,
		ImageLink:    nil,
		FollowIDs:    nil,
		Favorites:    nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now().Add(10),
	}

	userOut := rpc.UserOut{}.FromModel(user, "abc")

	bytes, err := json.Marshal(userOut)
	str := string(bytes)
	fmt.Println(err, str)
}
