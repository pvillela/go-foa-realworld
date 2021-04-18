package main

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
)

var (
	ErrXxx = errx.NewKind("xxx \"%v\"")
	ErrYyy = errx.NewKind("yyy \"%v\"")
	ErrZzz = errx.NewKind("zzz \"%v\"")
)

var bazErr errx.Errx

func baz() errx.Errx {
	bazErr = ErrXxx.Make(nil, "baz")
	return bazErr
}

func bar() errx.Errx {
	err := baz()
	err = ErrYyy.Make(err, "bar")
	return err
}

func foo() errx.Errx {
	err := bar()
	err = ErrZzz.Decorate(err, "foo")
	return err
}

type errW error

func main() {
	fmt.Println(ErrXxx.Make(nil))
	fmt.Println(ErrXxx.Make(nil, "foo"))
	fmt.Println(ErrXxx.Make(nil, "foo", "bar"))

	fooErr := foo()

	fmt.Println("\n---fmt.Printf(\"%+v\\n\", errW(fooErr))----------------------------------------------")
	fmt.Printf("%+v\n", errW(fooErr))

	fmt.Println("\n---fmt.Println(fooErr)----------------------------------------------")
	fmt.Println(fooErr)

	fmt.Println("\n---fmt.Println(error(fooErr))----------------------------------------------")
	fmt.Println(error(fooErr))

	fmt.Println("\n---fmt.Printf(\"%+v\\n\", fooErr)----------------------------------------------")
	fmt.Printf("%+v\n", fooErr)

	fmt.Println("\n---fmt.Printf(\"%+v\\n\", error(fooErr))----------------------------------------------")
	fmt.Printf("%+v\n", error(fooErr))

	fmt.Println("\n---fooErr.StackTrace())----------------------------------------------")
	fmt.Println(fooErr.StackTrace())

	fmt.Println("\n---errx.StackTrace(fooErr))----------------------------------------------")
	fmt.Println(errx.StackTraceOf(fooErr))

	fmt.Println("\n---fmt.Printf(\"%+v\\n\", errx.StackTrace(error(fooErr)))----------------------------------------------")
	fmt.Printf("%+v\n", errx.StackTraceOf(error(fooErr)))

	fmt.Println("\n---fmt.Printf(\"%+v\\n\", fooErr.DirectStackTrace())----------------------------------------------")
	fmt.Printf("%+v\n", fooErr.DirectStackTrace())

	fmt.Println("\n---fmt.Println(fooErr.Cause())----------------------------------------------")
	fmt.Println(fooErr.Cause())

	fmt.Println("\n---fmt.Println(fooErr.InnermostErrx())----------------------------------------------")
	fmt.Println(fooErr.InnermostErrx())

	fmt.Println("\n---fmt.Println(fooErr.InnermostCause())----------------------------------------------")
	fmt.Println(fooErr.InnermostCause())

	fmt.Println("\n---fooErr.CauseChain()----------------------------------------------")
	for _, err := range fooErr.CauseChain() {
		fmt.Println(err)
	}

	fmt.Println("\n---fooErr.ErrxChain()----------------------------------------------")
	for _, err := range fooErr.ErrxChain() {
		fmt.Println(err)
	}

	fmt.Println("\n===bazErr=====================================================================")

	fmt.Println("\n---fmt.Printf(\"%+v\\n\", errW(bazErr))----------------------------------------------")
	fmt.Printf("%+v\n", errW(bazErr))

	fmt.Println("\n---fmt.Println(bazErr)----------------------------------------------")
	fmt.Println(bazErr)

	fmt.Println("\n---fmt.Printf(\"%+v\\n\", bazErr)----------------------------------------------")
	fmt.Printf("%+v\n", bazErr)

	fmt.Println("\n---fmt.Println(errx.StackTrace(bazErr))----------------------------------------------")
	fmt.Println(errx.StackTraceOf(bazErr))

	fmt.Println("\n---fmt.Printf(\"%+v\\n\", errx.StackTrace(bazErr)----------------------------------------------")
	fmt.Printf("%+v\n", errx.StackTraceOf(bazErr))
}
