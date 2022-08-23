package web

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func WaitForHttpServer(baseUrl string, timeout time.Duration) error {
	startTime := time.Now()
	for {
		resp, err := http.Get(baseUrl)
		log.Debug("resp:", resp)
		log.Debug("err:", err)
		if err == nil {
			return nil
		}
		elapsed := time.Now().Sub(startTime)
		if elapsed > timeout {
			return fmt.Errorf("WaitForHttpServer timeout (%v) exceeded", timeout)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
