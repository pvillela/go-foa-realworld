package rpc

type RegisterIn struct {
	User struct {
		Username string
		Email    string
		Password string
	}
}
