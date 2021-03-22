package rpc

type UserAuthenticateIn struct {
	User struct {
		Email    string
		Password string
	}
}
