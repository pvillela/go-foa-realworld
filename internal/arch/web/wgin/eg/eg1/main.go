package main

import (
	"context"
	"fmt"
	"io"
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
	fmt.Println("authenticator ran\n", "tokenDetails:", tokenDetails)
	return true, tokenDetails.AccessToken, nil
}

var defaultReqCtxExtractor = web.DefaultReqCtxExtractor

var svcH = wgin.MakeStdFullBodySflHandler[In, Out](
	nil, false, defaultReqCtxExtractor, web.DefaultErrorHandler,
)(svc)

func main() {
	out := os.Stdout
	r, w := io.Pipe()
	mw := io.MultiWriter(out, w)
	bytes := []byte{99, 88}
	go func() {
		r.Read(bytes)
	}()
	mw.Write([]byte{1, 2})
	fmt.Println("bytes =", bytes)
	out.Write([]byte{3, 3, 3, 3, 3, 3, 3, 3, 3, 3})

	go func() {
		fmt.Println("11111111111111")
		router := gin.Default()
		fmt.Println("22222222222222")

		// Example for binding JSON ({"user": "manu", "password": "123"})

		router.POST("/loginJSON", svcH)
		router.POST("/loginJSON/:password", svcH)
		fmt.Println("333333333333333")

		// Listen and serve on 0.0.0.0:8080
		err := router.Run(":8080")
		fmt.Println("Server terminated:", err) // this never prints
	}()

	n := 5
	delta := 5 * time.Second
	for i := 1; i <= n; i++ {
		time.Sleep(delta)
		fmt.Println("Running server:", (time.Duration(i) * delta).Seconds(),
			"seconds of", (time.Duration(n) * delta).Seconds(), "seconds")
	}
	fmt.Println("read bytes =", bytes)
	fmt.Println("Exiting")
}
