package model

type User struct {
	Email    string
	Token    string
	Username string
	Bio      string
	Image    string // image link, nullable
}

type Author struct {
	Username  string
	Bio       string
	Image     string
	Following bool
}
