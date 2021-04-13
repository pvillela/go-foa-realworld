package main

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
)

var (
	ErrDuplicateKey = util.NewErrKind("ErrDuplicateKey", "duplicate key \"%v\"")
)

var err util.Err

func bar() {
	err = ErrDuplicateKey.MakeWithSt(true, "foo")
}

func foo() {
	bar()
}

type errW util.Err

func main() {
	fmt.Println(ErrDuplicateKey.MakeWithSt(false))
	fmt.Println(ErrDuplicateKey.MakeWithSt(false, "foo"))
	fmt.Println(ErrDuplicateKey.MakeWithSt(false, "foo", "bar"))
	foo()
	fmt.Printf("%+v\n", errW(err))
	fmt.Println("-------------------------------------------------")
	fmt.Println(err.StackTrace())
}
