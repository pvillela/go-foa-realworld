package rpc

type UpdateUserIn struct {
	User struct {
		Email string
		Bio   string
		Image string
	}
}
