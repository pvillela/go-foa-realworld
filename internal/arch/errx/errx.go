package errx

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
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
	StackTrace() StackTrace
	DirectStackTrace() StackTrace
	ErrxChain() []Errx
	CauseChain() []error
	InnermostCause() error
	InnermostErrx() Errx
}

// Interface verification
func _() {
	func(errx Errx) {}(&errxImpl{})
}

type StackTrace struct {
	inner errors.StackTrace
}

func (s StackTrace) Format(state fmt.State, verb rune) {
	s.inner.Format(state, verb)
}

/////////////////////
// Private types

type errxImpl struct {
	kind   *Kind
	args   []interface{}
	cause  error
	tracer stackTracer
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type dummyError struct{}

func (dummyError) Error() string { return "" }

/////////////////////
// Helper functions

func castToErrx(err error) *errxImpl {
	errx, ok := err.(*errxImpl)
	if ok {
		return errx
	}
	return nil
}

func stackTrace(err error) errors.StackTrace {
	tracer, ok := err.(stackTracer)
	if ok {
		return tracer.StackTrace()
	}
	return nil
}

func (e *errxImpl) msgWithArgs() string {
	return fmt.Sprintf(e.kind.msg, e.args...)
}

func (errx *errxImpl) traverseErrxChain(includeSelf bool, f func(*errxImpl) bool) {
	e := errx
	if !includeSelf {
		e = castToErrx(e.cause)
	}
	for e != nil {
		cont := f(e)
		if !cont {
			return
		}
		e = castToErrx(e.cause)
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

func (e *errxImpl) StackTrace() StackTrace {
	var trace StackTrace
	var cause error

	f := func(e *errxImpl) bool {
		if e.tracer != nil {
			trace = StackTrace{e.tracer.StackTrace()}
			return false
		}
		cause = e.cause
		return true
	}

	e.traverseErrxChain(true, f)

	if trace.inner != nil {
		return trace
	}

	// The innermost cause in the chain may be a stackTracer
	return StackTrace{stackTrace(cause)}
}

func (e *errxImpl) DirectStackTrace() StackTrace {
	if e.tracer != nil {
		return StackTrace{e.tracer.StackTrace()}
	}
	return StackTrace{}
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

/////////////////////
// For fmt.Printf support

func (e errxImpl) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			for _, errx := range e.ErrxChain() {
				_, _ = io.WriteString(s, errx.Msg())
				if trace := errx.DirectStackTrace(); trace.inner != nil {
					trace.Format(s, verb)
				}
				_, _ = io.WriteString(s, "\n")
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

/////////////////////
// Other public functions

func StackTraceOf(err error) StackTrace {
	switch e := err.(type) {
	case Errx:
		return e.StackTrace()
	case stackTracer:
		return StackTrace{e.StackTrace()}
	default:
		return StackTrace{}
	}
}
