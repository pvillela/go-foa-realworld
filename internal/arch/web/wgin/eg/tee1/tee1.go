package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	fn := logOutput()
	defer fn()
	var words = "Hello\nDoctor\nName\nContinue\nYesterday\nTomorrow"
	for i := 0; i < 10; i++ {
		log.Println(i)
		fmt.Println(i)
		fmt.Println(words)
	}
}

func logOutput() func() {
	//logfile := `logfile`
	// open file read/write | create if not exist | clear file at open if exists
	//f, _ := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	// save existing stdout | MultiWriter writes to saved stdout and file
	out := os.Stdout

	pr, pw := io.Pipe()

	go func() {
		newline := true
		bytes := make([]byte, 1000)
		for {
			n, _ := pr.Read(bytes)
			if newline {
				out.Write([]byte("pr="))
				newline = false
			}
			out.Write([]byte("begin****************************\n"))
			out.Write(bytes[:n])
			out.Write([]byte("end****************************\n"))
			if bytes[0] == '\n' {
				newline = true
			}
		}
	}()

	mw := io.MultiWriter(out, pw)

	// get pipe reader and writer | writes to pipe writer come out pipe reader
	r, w, _ := os.Pipe()

	// replace stdout,stderr with pipe writer | all writes to stdout, stderr will go through pipe instead (fmt.print, log)
	os.Stdout = w
	os.Stderr = w

	// writes with log.Print should also write to mw
	log.SetOutput(mw)

	//create channel to control exit | will block until all copies are finished
	exit := make(chan bool)

	go func() {
		// copy all reads from pipe to multiwriter, which writes to stdout and file
		_, _ = io.Copy(mw, r)
		// when r or w is closed copy will finish and true will be sent to channel
		exit <- true
	}()

	// function to be deferred in main until program exits
	return func() {
		// close writer then block on exit channel | this will let mw finish writing before the program exits
		_ = w.Close()
		<-exit
		// close file after all writes have finished
		_ = pw.Close()
	}

}
