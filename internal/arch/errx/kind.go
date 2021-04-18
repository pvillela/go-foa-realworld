package errx

import (
	"github.com/pkg/errors"
)

var NilKind = Kind{}

type Kind struct {
	msg *string // uses pointer to ensure uniqueness of each Kind instance even if with the same message
}

func KindOf(err error) Kind {
	errx, ok := err.(*errxImpl)
	if !ok {
		return NilKind
	}
	return errx.kind
}

func NewKind(msg string) Kind {
	str := new(string)
	*str = msg
	return Kind{msg: str}
}

func (s Kind) Make(cause error, args ...interface{}) Errx {
	err := errxImpl{kind: s}
	err.args = args
	err.cause = cause
	err.tracer = errors.WithStack(dummyError{}).(stackTracer)
	return &err
}

func (s Kind) Decorate(cause Errx, args ...interface{}) Errx {
	err := errxImpl{kind: s}
	err.args = args
	err.cause = cause
	return &err
}
