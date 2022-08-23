package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/arch/web/wgin"
	"github.com/pvillela/go-foa-realworld/internal/arch/web/wgin/eg"
	log "github.com/sirupsen/logrus"
	"net/http"
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

func processRequests() {
	callServer := func(path string, payload []byte) {
		fmt.Println()
		fmt.Println("payload:", string(payload))
		resp, err := http.Post("http://localhost:8080"+path, "application/json",
			bytes.NewBuffer(payload))
		errx.PanicOnError(err)

		var res map[string]any
		err = json.NewDecoder(resp.Body).Decode(&res)
		errx.PanicOnError(err)
		fmt.Println((resp.Status))
		fmt.Println("response:", res)
	}

	{
		path := "/loginJSON"
		payload := []byte(`{ "user": "manu", "password": "123" }`)
		callServer(path, payload)
	}

	{
		path := "/loginJSON?password=123"
		payload := []byte(`{ "user": "manu" }`)
		callServer(path, payload)
	}

	{
		path := "/loginJSON/123"
		payload := []byte(`{ "user": "manu" }`)
		callServer(path, payload)
	}

	{
		path := "/loginJSON"
		payload := []byte(`{ "user": "manu" }`)
		callServer(path, payload)
	}
}

func main() {
	// Create Gin engine
	router := gin.Default()

	// Define routes
	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.POST("/loginJSON", svcH)
	router.POST("/loginJSON/:password", svcH)

	// Launch Gin server
	go func() {
		// Listen and serve on 0.0.0.0:8080
		err := router.Run(":8080")

		// Lines below don't execute unless there is an error
		log.Fatal("Server terminated:", err)
	}()

	// Wait for server to be ready
	err := web.WaitForHttpServer("http://localhost:8080/", 100*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("***** Server ready")

	processRequests()

	fmt.Println("Exiting")
}
