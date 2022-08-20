package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/arch/web/wgin"
	"github.com/pvillela/go-foa-realworld/internal/arch/web/wgin/eg"
)

// Binding from JSON
type In struct {
	User string `form:"user" uri:"user" json:"user" binding:"required"`
	// Removed binding:"required" below to allow password to be passed via URL
	Password string `form:"password" uri:"password" json:"password"`
}

type Out struct {
	Out string
}

func svc(ctx context.Context, reqCtx web.RequestContext, in In) (Out, error) {
	username := reqCtx.Username
	fmt.Println("username = ", username)
	if in.User != "manu" || in.Password != "123" {
		return Out{}, fmt.Errorf("Invalid user='%v' or password='%v'", in.User, in.Password)
	}
	return Out{in.User + in.Password}, nil
}

func dummyAuthenticator(pReq *http.Request) (bool, *jwt.Token, error) {
	token, err := web.VerifiedJwtToken(pReq, eg.SecretKey)
	if err != nil {
		return false, token, err
	}
	fmt.Println("authenticator ran\n", "claims:", token)
	return true, token, nil
}

var authenticator = web.AuthenticatorC(eg.SecretKey)

var defaultReqCtxExtractor = web.DefaultReqCtxExtractor

var svcH = wgin.MakeStdFullBodySflHandler[In, Out](
	authenticator, true, defaultReqCtxExtractor, web.DefaultErrorHandler,
)(svc)

func main() {
	router := gin.Default()

	// Example for binding JSON ({"user": "manu", "password": "123"})

	router.POST("/loginJSON", svcH)
	router.POST("/loginJSON/:password", svcH)

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
