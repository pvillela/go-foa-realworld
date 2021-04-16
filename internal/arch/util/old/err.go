package util

import (
	"fmt"
	"runtime"
)

const maxStackTraceBytes = 10_000
const withStackTraceDefault = true

var NilErrKind = ErrKind{}

type ErrKind struct {
	msg *string // uses pointer to ensure uniqueness of each ErrKind instance even if with the same message
}

type Err struct {
	ErrKind
	args       []interface{}
	stackTrace string
	cause      error
}

func (e Err) Error() string {
	return fmt.Sprintf(*e.msg, e.args...)
}

func ErrKindOf(err error) ErrKind {
	myErr, ok := err.(Err)
	if !ok {
		return NilErrKind
	}
	return myErr.ErrKind
}

func NewErrKind(msg string) ErrKind {
	str := new(string)
	*str = msg
	return ErrKind{msg: str}
}

func (s ErrKind) Make(cause error, args ...interface{}) Err {
	return s.MakeWithSt(withStackTraceDefault, cause, args...)
}

func (s ErrKind) MakeWithSt(withStackTrace bool, cause error, args ...interface{}) Err {
	err := Err{ErrKind: s}
	err.args = args
	err.cause = cause
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
