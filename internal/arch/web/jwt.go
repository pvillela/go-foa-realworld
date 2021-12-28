package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/twinj/uuid"

	jwt "github.com/golang-jwt/jwt/v4"
)

type TokenDetails struct {
	AccessToken  *jwt.Token
	RefreshToken *jwt.Token
	AccessUuid   string
	RefreshUuid  string
}

func (s TokenDetails) AccessClaims() jwt.MapClaims {
	return s.AccessToken.Claims.(jwt.MapClaims)
}

func (s TokenDetails) Sub() string {
	return s.AccessClaims()["sub"].(string)
}

func (s TokenDetails) RefreshClaims() jwt.MapClaims {
	return s.RefreshToken.Claims.(jwt.MapClaims)
}

func (s TokenDetails) AccessExp() int64 {
	return s.AccessClaims()["exp"].(int64)
}

func (s TokenDetails) RefreshExp() int64 {
	return s.RefreshClaims()["exp"].(int64)
}

func CreateToken(
	userid string,
	secretKey []byte,
	accessMinutes int,
	refreshHours int,
) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AccessUuid = uuid.NewV4().String()
	td.RefreshUuid = uuid.NewV4().String()

	var err error

	// Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["sub"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * time.Duration(accessMinutes)).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	at.Raw, err = at.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	td.AccessToken = at

	// Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["sub"] = userid
	rtClaims["exp"] = time.Now().Add(time.Hour * time.Duration(refreshHours)).Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	rt.Raw, err = rt.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	td.RefreshToken = rt

	return td, nil
}

func extractToken(r *http.Request) string {
	// Typical header structure -- Authorization: Bearer <token>
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifiedJwtToken(r *http.Request, secretKey []byte) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
