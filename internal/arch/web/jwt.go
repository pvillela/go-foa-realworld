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
	AccessToken  string
	AtClaims     jwt.MapClaims
	RefreshToken string
	RtClaims     jwt.MapClaims
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func CreateToken(
	userid string,
	secretKey []byte,
	accessMinutes int,
	refreshHours int,
) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * time.Duration(accessMinutes)).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * time.Duration(refreshHours)).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["sub"] = userid
	atClaims["exp"] = td.AtExpires
	td.AtClaims = atClaims
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString(secretKey)
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["sub"] = userid
	rtClaims["exp"] = td.RtExpires
	td.RtClaims = rtClaims
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString(secretKey)
	if err != nil {
		return nil, err
	}
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
