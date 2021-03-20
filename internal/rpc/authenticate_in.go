package rpc

type AuthenticateIn struct {
	User struct {
		Email    string
		Password string
	}
}
