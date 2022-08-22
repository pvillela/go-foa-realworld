package wgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strings"
)

// GinLaunchAndSignal launches Gin on a given port in a separate goroutine and returns
// a channel that signals when the server is ready and a function to be deferred to
// close the pipe used in the implementation. The router parameter is a Gin Engine with
// configured routes.
func GinLaunchAndSignal(router *gin.Engine, port int) (serverReady chan bool, closePipe func()) {
	// Create memory pipe for tee with stdout
	pr, pw := io.Pipe()

	// Define tee
	mw := io.MultiWriter(os.Stdout, pw)

	// String in stdout that signals when server is ready
	serverReadyStr := "Listening and serving HTTP on"

	// Channel to signal when server is ready
	serverReady = make(chan bool)

	// Look for string in stdout and send server ready signal
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

	// Get os pipe reader and writer; writes to pipe writer come out pipe reader
	r, w, _ := os.Pipe()

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
		// Listen and serve on 0.0.0.0:port
		err := router.Run(fmt.Sprintf(":%v", port))
		fmt.Println("Server terminated:", err) // this never prints unless there is an error
	}()

	closePipe = func() {
		_ = w.Close()
		<-exit
	}

	return serverReady, closePipe
}
