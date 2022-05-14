package errx

import (
	"fmt"
)

/////////////////////
// Public types

// Errx defines an error type with support for error kinds, a cause chain, recursive error
// message, and a stack trace.
type Errx interface {
	error

	// Kind returns the error's Kind
	Kind() *Kind

	// Cause returns the error's cause, which may be nil.
	// This method is named for consistency with Go's standard errors package.
	Unwrap() error

	// Args returns the arguments that are substituted into KindMsg().
	Args() []any

	// KindMsg returns the raw message for the error's kind, i.e., Kind().msg.
	KindMsg() string

	// Msg returns the error's message with arguments substituted.
	Msg() string

	// RecursiveMsg returns a message string that combines the error messages of all errors
	// in the error's cause chain (which includes the error itself).
	RecursiveMsg() string

	// ErrxChain returns the error followed by all its preceding causes of type Errx.
	ErrxChain() []Errx

	// CauseChain returns the error followed by all its preceding causes.
	CauseChain() []error

	// InnermostCause returns the innermost cause in the error's cause chain.
	InnermostCause() error

	// InnermostErrx returns the innermost cause of type Errx in the error's cause chain.
	InnermostErrx() Errx

	// StackTrace returns a stack trace from the point where the error was created.
	StackTrace() string
}

// Interface verification
func _() {
	func(errx Errx) {}(&errxImpl{})
}

/////////////////////
// Private types

type errxImpl struct {
	kind                 *Kind
	args                 []any
	cause                error
	stack                []byte
	stackLinesToSuppress int
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

func (e *errxImpl) Error() string {
	return e.RecursiveMsg()
}

func (e *errxImpl) Kind() *Kind {
	return e.kind
}

func (e *errxImpl) Unwrap() error {
	return e.cause
}

func (e *errxImpl) Args() []any {
	return e.args
}

func (e *errxImpl) KindMsg() string {
	return e.kind.msg
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

// Guaranteed to return a non-nil result because *errxImpl can only be instantiated with
// the (*Kind).makeInternal constructor with sets a non-nil stack or with the (*Kind).Decorate
// constructor which requires an existing Errx as the cause which already must have a non-nil
// stack somewhere in its cause chain.
func (e *errxImpl) firstErrWithStack() *errxImpl {
	var err *errxImpl
	f := func(e *errxImpl) bool {
		if e.stack != nil {
			err = e
			return false
		}
		return true
	}

	e.traverseErrxChain(true, f)

	return err
}

// StackTrace returns a string that contains both the
// error message and the callstack.
func (e *errxImpl) StackTrace() string {
	ews := e.firstErrWithStack() // guaranteed to be non-nil
	cutoffLineIndex := 0
	newlineCount := 0
	for i, b := range ews.stack {
		cutoffLineIndex = i + 1
		if b == '\n' {
			newlineCount++
		}
		if newlineCount == ews.stackLinesToSuppress*2+1 {
			break
		}
	}
	trimmedStack := ews.stack[cutoffLineIndex:]
	return "errx.Errx: " + e.Error() + "\n" + string(trimmedStack)
}

/////////////////////
// Factory functions. See also kind.go for additional factory functions.

// Helper method to create an Errx whose Kind is defined on-the-fly using msg.
func newErrxInternal(cause error, msg string, stackLinesToSuppress int) Errx {
	kind := NewKind(msg)
	err := kind.makeInternal(cause, stackLinesToSuppress)
	return err
}

// NewErrx creates an Errx whose Kind is defined on-the-fly using msg.
func NewErrx(cause error, msg string) Errx {
	return newErrxInternal(cause, msg, 4)
}

// ErrxOf creates an Errx from r.
// If r is nil, nil is returned.
// If r is an Errx, r is returned.
// If r is an error, NewErrx is used to instantiate an Errx with r as its cause.
// Otherwise, NewErrx is used to instantiate an Errx with nil as the cause argument
// and r's string rendering as the msg argument.
func ErrxOf(r any) Errx {
	if r == nil {
		return nil
	}
	var err error
	switch r.(type) {
	case error:
		err = r.(error)
	default:
		err = nil
	}
	errX, ok := err.(Errx)
	if !ok {
		if err != nil {
			errX = newErrxInternal(err, ".", 4)
		} else {
			errX = newErrxInternal(nil, fmt.Sprintf("%v", r), 4)
		}
	}
	return errX
}

/////////////////////
// Other public functions

// StackTraceOf returns err.StackTrace() if err is of type Errx, an empty string otherwise.
func StackTraceOf(err error) string {
	switch e := err.(type) {
	case Errx:
		return e.StackTrace()
	default:
		return ""
	}
}
