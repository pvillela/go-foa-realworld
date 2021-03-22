package rpc

type UserRegisterIn struct {
	User struct {
		Username string
		Email    string
		Password string
	}
}
