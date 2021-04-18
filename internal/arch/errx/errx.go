package errx

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
)

/////////////////////
// Types

type Errx interface {
	error
	Kind() Kind
	Cause() error
	Args() []interface{}
	FullMsg() string
	DirectMsg() string
	StackTrace() errors.StackTrace
	DirectStackTrace() errors.StackTrace
	ErrxChain() []Errx
	CauseChain() []error
	InnermostCause() error
	InnermostErrx() Errx
}

// Interface verification
func _() {
	func(errx Errx) {}(&errxImpl{})
}

type errxImpl struct {
	kind            Kind
	args            []interface{}
	cause           error
	errorsCompanion error // for use of errors.New and errors.Wrap
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
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

func stackTrace(err error) errors.StackTrace {
	tracer, ok := err.(stackTracer)
	if ok {
		return tracer.StackTrace()
	}
	return nil
}

func (e *errxImpl) msgWithArgs() string {
	return fmt.Sprintf(*e.kind.msg, e.args...)
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
	return e.errorsCompanion.Error()
}

func (e *errxImpl) Kind() Kind {
	return e.kind
}

func (e *errxImpl) Cause() error {
	return e.cause
}

func (e *errxImpl) Args() []interface{} {
	return e.args
}

func (e *errxImpl) FullMsg() string {
	return e.Error()
}

func (e *errxImpl) DirectMsg() string {
	return e.msgWithArgs()
}

func (e *errxImpl) StackTrace() errors.StackTrace {
	var trace errors.StackTrace
	var cause error

	f := func(e *errxImpl) bool {
		if trace = stackTrace(e.errorsCompanion); trace != nil {
			return false
		}
		cause = e.cause
		return true
	}

	e.traverseErrxChain(true, f)

	if trace != nil {
		return trace
	}

	// The innermost cause in the chain may be a stackTracer
	return stackTrace(cause)
}

func (e *errxImpl) DirectStackTrace() errors.StackTrace {
	tracer, ok := e.errorsCompanion.(stackTracer)
	if ok {
		return tracer.StackTrace()
	}
	return nil
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
				_, _ = io.WriteString(s, errx.DirectMsg())
				if trace := errx.DirectStackTrace(); trace != nil {
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

func StackTrace(err error) errors.StackTrace {
	tracer, ok := err.(stackTracer)
	if ok {
		return tracer.StackTrace()
	}
	return nil
}
