package wgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"io"
	"net/http"
	"os"
	//"strings"
	"time"
)

// GinLaunchAndSignal launches Gin on a given port in a separate goroutine and returns
// a channel that signals when the server is ready and a function to be deferred to
// close the pipe used in the implementation. The router parameter is a Gin Engine with
// configured routes.
func GinLaunchAndSignal(
	router *gin.Engine,
	port int,
) (serverReady chan bool, closePipe func(), err error) {
	// Create memory pipe for tee with stdout
	pr, pw := io.Pipe()
	util.Ignore(pr)
	// Define tee
	mw := io.MultiWriter(os.Stdout, pw)

	// String in stdout that signals when server is ready
	//serverReadyStr := "Listening and serving HTTP on"

	// Channel to signal when server is ready
	serverReady = make(chan bool)

	// Look for string in stdout and send server ready signal
	//go func() {
	//	const bufSize = 100
	//	bytes := make([]byte, bufSize)
	//	sb := strings.Builder{}
	//	found := false
	//	for {
	//		n, _ := pr.Read(bytes)
	//		if !found {
	//			sb.Write(bytes[:n])
	//			if strings.Contains(sb.String(), serverReadyStr) {
	//				found = true
	//				//// wait because Gin outputs serverReadyStr before it is really ready
	//				//time.Sleep(100 * time.Millisecond)
	//				serverReady <- true
	//			}
	//		}
	//	}
	//}()

	go func() {
		for {
			resp, err := http.Get("http://localhost:8080/")
			fmt.Println("resp:", resp)
			fmt.Println("err:", err)
			if err == nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		serverReady <- true
	}()

	// Get os pipe reader and writer; writes to pipe writer come out pipe reader
	r, w, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	// Create channel to control exit; will block until all copies are finished
	exit := make(chan bool)

	go func() {
		// copy all reads from os pipe to multiwriter, which writes to stdout and memory pipe
		_, _ = io.Copy(mw, r)
		// when r or w is closed copy will finish and true will be sent to channel
		exit <- true
	}()

	// Redefine Gin writer to use tee
	gin.DefaultWriter = w

	go func() {
		time.Sleep(300 * time.Millisecond)

		// Listen and serve on 0.0.0.0:port
		err := router.Run(fmt.Sprintf(":%v", port))

		// Lines below don't execute unless there is an error
		fmt.Println("Server terminated:", err)
		os.Exit(1)
	}()

	closePipe = func() {
		_ = w.Close()
		<-exit
	}

	return serverReady, closePipe, nil
}
