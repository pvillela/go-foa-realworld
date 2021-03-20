package model

type Profile struct {
	Username  string
	Bio       string
	Image     string // image link, nullable
	Following bool
}
