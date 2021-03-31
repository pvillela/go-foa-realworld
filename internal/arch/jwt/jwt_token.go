package jwt

import (
	"errors"
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"crypto/sha256"
)

var tokenTimeToLive = time.Hour * 2

func Hash(salt []byte, text string) string {
	h := sha256.New()
	bytes := make([]byte, len(salt)+len(text))
	bytes = append(bytes, salt...)
	bytes = append(bytes, []byte(text)...)
	h.Write(bytes)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// UserGenToken generates a JWT token for a user
func UserGenToken(user *model.User) (string, error) {
	if user.Name == "" {
		return "", errors.New("can't generate token for empty user")
	}

	return jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, newUserClaims(user.Name, tokenTimeToLive)).
		SignedString(user.PasswordSalt)
}

func UserGetNameFromToken(user model.User, tokenStr string) (string, error) {
	token, err := jwtgo.Parse(
		tokenStr,
		func(token *jwtgo.Token) (interface{}, error) {
			// TODO: this function doesn't check the token signature
			return "", nil
		},
	)
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwtgo.StandardClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", errors.New("problem with jwt token")
}

// newUserClaims : constructor of userClaims
func newUserClaims(username string, ttl time.Duration) *jwtgo.StandardClaims {
	return &jwtgo.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Issuer:    "real-world-demo-backend",
	}
}
