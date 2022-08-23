package main

import (
	"context"
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"net/http"
	"os"
	"time"

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
	tokenDetails, _ := web.CreateToken("me", eg.SecretKey, 100_000_000, 2_000_000)
	fmt.Println("authenticator ran\n", "token:", tokenDetails.AccessToken.Raw)
	return true, tokenDetails.AccessToken, nil
}

var defaultReqCtxExtractor = web.DefaultReqCtxExtractor

var svcH = wgin.MakeStdFullBodySflHandler[In, Out](
	dummyAuthenticator, true, defaultReqCtxExtractor, web.DefaultErrorHandler,
)(svc)

func main() {
	router := gin.Default()

	// Example for binding JSON ({"user": "manu", "password": "123"})

	router.POST("/loginJSON", svcH)
	router.POST("/loginJSON/:password", svcH)

	// Launch the server
	go func() {
		// Listen and serve on 0.0.0.0:8080
		err := router.Run(":8080")
		errx.PanicOnError(err)
	}()

	// Wait for server to be ready
	err := web.WaitForHttpServer("http://localhost:8080/", 100*time.Millisecond)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("***** Server ready")

	// Keep server running for n * delta seconds

	n := 5
	delta := 5 * time.Second
	for i := 1; i <= n; i++ {
		time.Sleep(delta)
		fmt.Println("Running server:", (time.Duration(i) * delta).Seconds(),
			"seconds of", (time.Duration(n) * delta).Seconds(), "seconds")
	}
	fmt.Println("Exiting")
}
