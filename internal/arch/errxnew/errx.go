package errxnew

import (
	"bytes"
	"fmt"
	"github.com/go-errors/errors"
)

/////////////////////
// Public types

type Errx interface {
	error
	Kind() *Kind
	Cause() error
	Args() []interface{}
	Msg() string
	RecursiveMsg() string
	ErrxChain() []Errx
	CauseChain() []error
	InnermostCause() error
	InnermostErrx() Errx
	StackTrace() string
}

// Interface verification
func _() {
	func(errx Errx) {}(&errxImpl{})
}

/////////////////////
// Private types

type errxImpl struct {
	err   *errors.Error
	kind  *Kind
	args  []interface{}
	cause error
}

type stackTracer interface {
	StackTrace() string
}

/////////////////////
// Helper functions

func castToErrx(err error) *errxImpl {
	errx, ok := err.(*errxImpl)
	if ok {
		return errx
	}
	return nil
}

func (e *errxImpl) msgWithArgs() string {
	return fmt.Sprintf(e.kind.msg, e.args...)
}

func (e *errxImpl) traverseErrxChain(includeSelf bool, f func(*errxImpl) bool) {
	err := e
	if !includeSelf {
		err = castToErrx(err.cause)
	}
	for err != nil {
		cont := f(err)
		if !cont {
			return
		}
		err = castToErrx(err.cause)
	}
	return
}

/////////////////////
// Methods

func (e *errxImpl) Err() *errors.Error {
	return e.err
}

func (e *errxImpl) Error() string {
	return e.RecursiveMsg()
}

func (e *errxImpl) Kind() *Kind {
	return e.kind
}

func (e *errxImpl) Cause() error {
	return e.cause
}

func (e *errxImpl) Args() []interface{} {
	return e.args
}

func (e *errxImpl) Msg() string {
	return e.msgWithArgs()
}

func (e *errxImpl) RecursiveMsg() string {
	str := e.Msg()
	if cause := e.cause; cause != nil {
		str = str + " -> " + cause.Error()
	}
	return str
}

func (e *errxImpl) ErrxChain() []Errx {
	chain := make([]Errx, 0, 1)
	f := func(e *errxImpl) bool {
		chain = append(chain, e)
		return true
	}
	e.traverseErrxChain(true, f)
	return chain
}

func (e *errxImpl) CauseChain() []error {
	chain := make([]error, 0, 1)
	f := func(e *errxImpl) bool {
		chain = append(chain, e.cause)
		return true
	}
	e.traverseErrxChain(true, f)
	return chain
}

func (e *errxImpl) InnermostErrx() Errx {
	var innermost Errx
	f := func(e *errxImpl) bool {
		innermost = e
		return true
	}
	e.traverseErrxChain(true, f)
	return innermost
}

func (e *errxImpl) InnermostCause() error {
	var cause error
	f := func(e *errxImpl) bool {
		cause = e.cause
		return true
	}
	e.traverseErrxChain(true, f)
	return cause
}

// Stack returns the callstack formatted the same way that go does
// in runtime/debug.Stack()
func (e *errxImpl) stack() []byte {
	buf := bytes.Buffer{}

	for _, frame := range e.err.StackFrames()[1:] {
		buf.WriteString(frame.String())
	}

	return buf.Bytes()
}

func (e *errxImpl) StackTrace() string {
	stack := e.stack()
	return "errx.Errx: " + e.Error() + "\n" + string(stack)
}

/////////////////////
// Other public functions

func StackTraceOf(err error) string {
	switch e := err.(type) {
	case Errx:
		return e.StackTrace()
	case stackTracer:
		return e.StackTrace()
	default:
		return ""
	}
}
