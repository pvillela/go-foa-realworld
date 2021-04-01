package rpc

type UserUpdateIn struct {
	User struct {
		Email string
		Bio   *string
		Image string
	}
}
