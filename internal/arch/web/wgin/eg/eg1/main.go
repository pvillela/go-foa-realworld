package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/arch/web/wgin"
	"github.com/pvillela/go-foa-realworld/internal/arch/web/wgin/eg"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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
	// save existing stdout | MultiWriter writes to saved stdout and pipe
	out := os.Stdout

	pr, pw := io.Pipe()

	serverReadyStr := "Listening and serving HTTP on"
	serverReady := make(chan bool)

	go func() {
		const bufSize = 100
		bytes := make([]byte, bufSize)
		sb := strings.Builder{}
		found := false
		for {
			n, _ := pr.Read(bytes)
			if !found {
				sb.Write(bytes[:n])
				if strings.Contains(sb.String(), serverReadyStr) {
					found = true
					serverReady <- true
				}
			}
		}
	}()

	mw := io.MultiWriter(out, pw)

	// get pipe reader and writer | writes to pipe writer come out pipe reader
	r, w, _ := os.Pipe()

	// replace stdout,stderr with pipe writer | all writes to stdout, stderr will go through pipe instead (fmt.print, log)
	os.Stdout = w
	os.Stderr = w

	//create channel to control exit | will block until all copies are finished
	exit := make(chan bool)

	go func() {
		// copy all reads from pipe to multiwriter, which writes to stdout and file
		_, _ = io.Copy(mw, r)
		// when r or w is closed copy will finish and true will be sent to channel
		exit <- true
	}()

	go func() {
		gin.DefaultWriter = w

		router := gin.Default()

		// Example for binding JSON ({"user": "manu", "password": "123"})

		router.POST("/loginJSON", svcH)
		router.POST("/loginJSON/:password", svcH)

		// Listen and serve on 0.0.0.0:8080
		err := router.Run(":8080")
		fmt.Println("Server terminated:", err) // this never prints unless there is an error
	}()

	<-serverReady
	//out.Write([]byte("***** Server is ready\n"))
	fmt.Println("***** Server is ready")
	n := 5
	delta := 5 * time.Second
	for i := 1; i <= n; i++ {
		time.Sleep(delta)
		fmt.Println("Running server:", (time.Duration(i) * delta).Seconds(),
			"seconds of", (time.Duration(n) * delta).Seconds(), "seconds")
	}
	fmt.Println("Exiting")
	w.Close()
	<-exit
}
