package main

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/olderr"
)

var (
	ErrDuplicateKey   = util.NewErrKind("database %v - duplicate key \"%v\"")
	ErrRecordNotFound = util.NewErrKind("database %v - record not found with key \"%v\"")
)

func bar() util.Err {
	return ErrDuplicateKey.Make(nil, "mydb", "bar")
}

func foo() util.Err {
	err := bar()
	return ErrRecordNotFound.Make(err, "mydb", "foo")
}

type errW util.Err

func main() {
	fmt.Println(ErrDuplicateKey.Make(nil, "mydb"))
	fmt.Println(ErrDuplicateKey.Make(nil, "mydb", "foo"))
	fmt.Println(ErrDuplicateKey.Make(nil, "mydb", "foo", "bar"))
	err := foo()
	fmt.Printf("*** Printf errW -> %+v\n", errW(err))
	fmt.Println("*** Println err -> ", err)
	fmt.Printf("*** Printf err -> %+v\n", err)
	fmt.Println("---err.StackTrace()----------------------------------------------")
	fmt.Println(err.StackTrace())
	fmt.Println("---err.Cause().(util.Err).StackTrace()----------------------------------------------")
	fmt.Println(err.Cause().(util.Err).StackTrace())
	fmt.Println("---err.StackTrace()----------------------------------------------")
	fmt.Printf("%+v\n", err.StackTrace())
}
