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
	err = ErrDuplicateKey.Make(true, "foo")
}

func foo() {
	bar()
}

type errW util.Err

func main() {
	fmt.Println(ErrDuplicateKey.Make(false))
	fmt.Println(ErrDuplicateKey.Make(false, "foo"))
	fmt.Println(ErrDuplicateKey.Make(false, "foo", "bar"))
	foo()
	fmt.Printf("%+v\n", errW(err))
	fmt.Println("-------------------------------------------------")
	fmt.Println(err.StackTrace())
}
