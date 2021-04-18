package main

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
)

var (
	ErrDuplicateKey   = errx.NewKind("database %v - duplicate key \"%v\"")
	ErrRecordNotFound = errx.NewKind("database %v - record not found with key \"%v\"")
	ErrXxx            = errx.NewKind("xxx \"%v\"")
)

var bazErr errx.Errx

func baz() errx.Errx {
	bazErr = ErrDuplicateKey.Make(nil, "mydb", "baz")
	return bazErr
}

func bar() errx.Errx {
	err := baz()
	err = ErrRecordNotFound.Make(err, "mydb", "bar")
	return err
}

func foo() errx.Errx {
	err := bar()
	err = ErrXxx.Decorate(err, "foo")
	return err
}

type errW error

func main() {
	fmt.Println(ErrDuplicateKey.Make(nil, "mydb"))
	fmt.Println(ErrDuplicateKey.Make(nil, "mydb", "foo"))
	fmt.Println(ErrDuplicateKey.Make(nil, "mydb", "foo", "bar"))

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
	fmt.Println(errx.StackTrace(fooErr))

	fmt.Println("\n---fmt.Printf(\"%+v\\n\", errx.StackTrace(error(fooErr)))----------------------------------------------")
	fmt.Printf("%+v\n", errx.StackTrace(error(fooErr)))

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
	fmt.Println(errx.StackTrace(bazErr))

	fmt.Println("\n---fmt.Printf(\"%+v\\n\", errx.StackTrace(bazErr)----------------------------------------------")
	fmt.Printf("%+v\n", errx.StackTrace(bazErr))
}
