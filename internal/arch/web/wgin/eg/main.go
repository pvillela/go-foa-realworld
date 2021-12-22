package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvillela/go-foa-realworld/internal/arch/web"
	"github.com/pvillela/go-foa-realworld/internal/arch/web/wgin"
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

func svc(_ string, in In) (Out, error) {
	if in.User != "manu" || in.Password != "123" {
		return Out{}, fmt.Errorf("Invalid user='%v' or password='%v'", in.User, in.Password)
	}
	return Out{in.User + in.Password}, nil
}

func authenticator(pReq *http.Request) error {
	// TODO: implement this
	fmt.Println("authenticator ran")
	return nil
}

var defaultReqCtxExtractor = web.MakeDefaultReqCtxExtractor([]byte("1234567890"))

var svcH = wgin.MakeStdFullBodySflToHandler[In, Out](
	authenticator, defaultReqCtxExtractor, web.DefaultErrorHandler,
)(svc)

func main() {
	router := gin.Default()

	// Example for binding JSON ({"user": "manu", "password": "123"})

	router.POST("/loginJSON", svcH)
	router.POST("/loginJSON/:password", svcH)

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
