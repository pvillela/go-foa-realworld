package main

import (
	"fmt"
	"github.com/pvillela/go-foa-realworld/internal/arch/errxnewbasic"
)

var (
	ErrXxx = errxnewbasic.NewKind("xxx \"%v\"")
	ErrYyy = errxnewbasic.NewKind("yyy \"%v\"", ErrXxx, ErrXxx)
	ErrZzz = errxnewbasic.NewKind("zzz \"%v\"", ErrYyy)
	ErrWww = errxnewbasic.NewKind("www \"%v\"", ErrYyy, ErrZzz)
)

var bazErr errxnewbasic.Errx

func baz() errxnewbasic.Errx {
	bazErr = ErrXxx.Make(nil, "baz")
	return bazErr
}

func bar() errxnewbasic.Errx {
	err := baz()
	err = ErrYyy.Make(err, "bar")
	return err
}

func foo() errxnewbasic.Errx {
	err := bar()
	err = ErrZzz.Make(err, "foo")
	return err
}

type errW error

func main() {
	fmt.Println(ErrXxx.Make(nil))
	fmt.Println(ErrXxx.Make(nil, "foo"))
	fmt.Println(ErrXxx.Make(nil, "foo", "bar"))

	fooErr := foo()

	fmt.Println("\n---errW(fooErr)--------------------------------------------")
	fmt.Printf("%+v\n", errW(fooErr))

	fmt.Println("\n---fmt.Println(fooErr)----------------------------------------------")
	fmt.Println(fooErr)

	fmt.Println("\n---fmt.Println(error(fooErr)--------------------------------------------")
	fmt.Println(error(fooErr))

	fmt.Println("\n---error(fooErr)--------------------------------------------")
	fmt.Printf("%+v\n", error(fooErr))

	fmt.Println("\n---fooErr.StackTrace()--------------------------------------------")
	fmt.Println(fooErr.StackTrace())

	fmt.Println("\n---errxnewbasic.StackTrace(fooErr)--------------------------------------------")
	fmt.Println(errxnewbasic.StackTraceOf(fooErr))

	fmt.Println("\n---errxnewbasic.StackTrace(error(fooErr))--------------------------------------------")
	fmt.Printf("%+v\n", errxnewbasic.StackTraceOf(error(fooErr)))

	fmt.Println("\n---fmt.Println(fooErr.Cause()--------------------------------------------")
	fmt.Println(fooErr.Cause())

	fmt.Println("\n---fmt.Println(fooErr.InnermostErrx()--------------------------------------------")
	fmt.Println(fooErr.InnermostErrx())

	fmt.Println("\n---fmt.Println(fooErr.InnermostCause()--------------------------------------------")
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

	fmt.Println("\n---errW(bazErr)--------------------------------------------")
	fmt.Printf("%+v\n", errW(bazErr))

	fmt.Println("\n---fmt.Println(bazErr)----------------------------------------------")
	fmt.Println(bazErr)

	fmt.Println("\n---fmt.Println(errxnewbasic.StackTrace(bazErr)--------------------------------------------")
	fmt.Println(errxnewbasic.StackTraceOf(bazErr))

	fmt.Println("\n---errxnewbasic.StackTrace(bazErr)----------------------------------------------")
	fmt.Printf("%+v\n", errxnewbasic.StackTraceOf(bazErr))

	fmt.Println("\n===SubKinds=====================================================================")
	fmt.Println()

	deref := func(m map[*errxnewbasic.Kind]struct{}) []errxnewbasic.Kind {
		slice := make([]errxnewbasic.Kind, 0, len(m))
		for kind, _ := range m {
			slice = append(slice, *kind)
		}
		return slice
	}

	fmt.Println("ErrXxx.SuperKinds()", deref(ErrXxx.SuperKinds()))
	fmt.Println("ErrXxx.IsSubKindOf(ErrXxx)", ErrXxx.IsSubKindOf(ErrXxx))
	fmt.Println("ErrXxx.IsSubKindOf(ErrYyy)", ErrXxx.IsSubKindOf(ErrYyy))
	fmt.Println()

	fmt.Println("ErrYyy.SuperKinds()", deref(ErrYyy.SuperKinds()))
	fmt.Println("ErrYyy.IsSubKindOf(ErrXxx)", ErrYyy.IsSubKindOf(ErrXxx))
	fmt.Println("ErrYyy.IsSubKindOf(ErrZzz)", ErrYyy.IsSubKindOf(ErrZzz))
	fmt.Println()

	fmt.Println("ErrZzz.SuperKinds()", deref(ErrZzz.SuperKinds()))
	fmt.Println("ErrZzz.IsSubKindOf(ErrXxx)", ErrZzz.IsSubKindOf(ErrXxx))
	fmt.Println("ErrZzz.IsSubKindOf(ErrYyy)", ErrZzz.IsSubKindOf(ErrYyy))
	fmt.Println("ErrZzz.IsSubKindOf(ErrZzz)", ErrZzz.IsSubKindOf(ErrZzz))
	fmt.Println("ErrZzz.IsSubKindOf(ErrWww)", ErrZzz.IsSubKindOf(ErrWww))
	fmt.Println()

	fmt.Println("ErrWww.SuperKinds()", deref(ErrWww.SuperKinds()))
	fmt.Println("ErrWww.IsSubKindOf(ErrXxx)", ErrWww.IsSubKindOf(ErrXxx))
	fmt.Println("ErrWww.IsSubKindOf(ErrYyy)", ErrWww.IsSubKindOf(ErrYyy))
	fmt.Println("ErrWww.IsSubKindOf(ErrZzz)", ErrWww.IsSubKindOf(ErrZzz))
	fmt.Println("ErrWww.IsSubKindOf(ErrWww)", ErrWww.IsSubKindOf(ErrWww))
}
