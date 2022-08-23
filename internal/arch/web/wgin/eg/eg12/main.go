package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/arch/web/wgin"
	"github.com/pvillela/go-foa-realworld/internal/arch/web/wgin/eg"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

// Binding from JSON
type In struct {
	User string `form:"user" uri:"user" json:"user" binding:"required"`
	// Removed binding:"required" below to allow password to be passed via URL
	Password string `form:"password" uri:"password" json:"password"`
}

type Out1 struct {
	UserId   string
	Password string
	Token    string
}

func svc1(ctx context.Context, reqCtx web.RequestContext, in In) (Out1, error) {
	if in.User != "manu" || in.Password != "123" {
		return Out1{}, fmt.Errorf("Invalid user='%v' or password='%v'", in.User, in.Password)
	}
	tokenDetails, _ := web.CreateToken("me", eg.SecretKey, 100_000_000, 2_000_000)
	return Out1{UserId: in.User, Password: in.Password, Token: tokenDetails.AccessToken.Raw}, nil
}

type Out2 struct {
	Out string
}

func svc2(ctx context.Context, reqCtx web.RequestContext, in In) (Out2, error) {
	username := reqCtx.Username
	fmt.Println("username = ", username)
	if in.User != "manu" || in.Password != "123" {
		return Out2{}, fmt.Errorf("Invalid user='%v' or password='%v'", in.User, in.Password)
	}
	return Out2{in.User + in.Password}, nil
}

var authenticator = web.AuthenticatorC(eg.SecretKey)

var defaultReqCtxExtractor = web.DefaultReqCtxExtractor

var svc1H = wgin.MakeStdFullBodySflHandler[In, Out1](
	nil, false, defaultReqCtxExtractor, web.DefaultErrorHandler,
)(svc1)

var svc2H = wgin.MakeStdFullBodySflHandler[In, Out2](
	authenticator, true, defaultReqCtxExtractor, web.DefaultErrorHandler,
)(svc2)

func PostWithHeaders(
	url, contentType string,
	body io.Reader,
	headers map[string]string,
) (resp *http.Response, err error) {
	client := http.DefaultClient
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return client.Do(req)
}

func processRequests() {
	callServerWithHeaders := func(
		path string,
		payload []byte,
		headers map[string]string,
	) (res map[string]any) {
		fmt.Println("payload:", string(payload))

		resp, err := PostWithHeaders("http://localhost:8080"+path, "application/json",
			bytes.NewBuffer(payload), headers)
		errx.PanicOnError(err)

		err = json.NewDecoder(resp.Body).Decode(&res)
		errx.PanicOnError(err)
		fmt.Println((resp.Status))
		fmt.Println("response:", res)
		fmt.Println()

		return res
	}

	callServer := func(path string, payload []byte) (res map[string]any) {
		return callServerWithHeaders(path, payload, nil)
	}

	var res map[string]any
	{
		fmt.Println("=== Good request -- all in payload")
		path := "/loginJSON1"
		payload := []byte(`{ "user": "manu", "password": "123" }`)
		res = callServer(path, payload)
	}
	authHeaderMap := map[string]string{
		"Authorization": "Bearer " + res["Token"].(string),
	}

	{
		fmt.Println("=== Good request -- password in query parameter")
		path := "/loginJSON1?password=123"
		payload := []byte(`{ "user": "manu" }`)
		callServer(path, payload)
	}

	{
		fmt.Println("=== Good request -- password in path")
		path := "/loginJSON1/123"
		payload := []byte(`{ "user": "manu" }`)
		callServer(path, payload)
	}

	{
		fmt.Println("=== Bad request -- missing password")
		path := "/loginJSON1"
		payload := []byte(`{ "user": "manu" }`)
		callServer(path, payload)
	}

	{
		fmt.Println("=== Good request with token -- all in payload")
		path := "/loginJSON2"
		payload := []byte(`{ "user": "manu", "password": "123" }`)
		res = callServerWithHeaders(path, payload, authHeaderMap)
	}

	{
		fmt.Println("=== Good request with token -- some data in query parameter")
		path := "/loginJSON2?password=123"
		payload := []byte(`{ "user": "manu" }`)
		callServerWithHeaders(path, payload, authHeaderMap)
	}

	{
		fmt.Println("=== Good request with token -- some data in path")
		path := "/loginJSON2/123"
		payload := []byte(`{ "user": "manu" }`)
		callServerWithHeaders(path, payload, authHeaderMap)
	}

	{
		fmt.Println("=== Bad request with token -- missing required input")
		path := "/loginJSON2"
		payload := []byte(`{ "user": "manu" }`)
		callServerWithHeaders(path, payload, authHeaderMap)
	}

	{
		fmt.Println("=== Bad request with token -- missing token")
		path := "/loginJSON2"
		payload := []byte(`{ "user": "manu", "password": "123" }`)
		callServerWithHeaders(path, payload, nil)
	}

	{
		fmt.Println("=== Bad request with token -- expired token")
		path := "/loginJSON2"
		payload := []byte(`{ "user": "manu", "password": "123" }`)
		authHeaderMap := map[string]string{
			"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjE2Mjg2NTgyLWM1NDEtNDM2NS05MGJiLTM5NTg0MDE4NTczNyIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTY0MDI1MjQ0NSwic3ViIjoibWUifQ.zRBG9My6-QKJ3-J7-5wX6GzoPhpFpKctbwnfQDAAuB4",
		}
		callServerWithHeaders(path, payload, authHeaderMap)
	}

	{
		fmt.Println("=== Bad request with token -- token with invalid signature")
		path := "/loginJSON2"
		payload := []byte(`{ "user": "manu", "password": "123" }`)
		authHeaderMap := map[string]string{
			"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjE2Mjg2NTgyLWM1NDEtNDM2NS05MGJiLTM5NTg0MDE4NTczNyIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTY0MDI1MjQ0NSwic3ViIjoibWUifQ.ZRBG9My6-QKJ3-J7-5wX6GzoPhpFpKctbwnfQDAAuB4",
		}
		callServerWithHeaders(path, payload, authHeaderMap)
	}
}

func main() {
	// Create Gin engine
	router := gin.Default()

	// Define routes
	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.POST("/loginJSON1", svc1H)
	router.POST("/loginJSON1/:password", svc1H)
	router.POST("/loginJSON2", svc2H)
	router.POST("/loginJSON2/:password", svc2H)

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
	fmt.Println()

	processRequests()

	fmt.Println("Exiting")
}
