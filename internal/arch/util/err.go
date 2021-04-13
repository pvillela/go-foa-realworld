package util

import (
	"fmt"
	"runtime"
)

const maxStackTraceBytes = 10_000
const withStackTraceDefault = true

var NilErrKind = ErrKind{}

type ErrKind struct {
	msg string
}

type Err struct {
	ErrKind
	args       []interface{}
	stackTrace string
}

func (e Err) Error() string {
	return fmt.Sprintf(e.msg, e.args...)
}

func ErrKindOf(err error) ErrKind {
	myErr, ok := err.(Err)
	if !ok {
		return NilErrKind
	}
	return myErr.ErrKind
}

func NewErrKind(msg string) ErrKind {
	return ErrKind{msg: msg}
}

func (s ErrKind) Make(args ...interface{}) Err {
	return s.MakeWithSt(withStackTraceDefault, args...)
}

func (s ErrKind) MakeWithSt(withStackTrace bool, args ...interface{}) Err {
	err := Err{ErrKind: s}
	err.args = args
	if withStackTrace {
		st := make([]byte, maxStackTraceBytes)
		runtime.Stack(st, false)
		err.stackTrace = string(st)
	}
	return err
}

func (e Err) StackTrace() string {
	return e.stackTrace
}
