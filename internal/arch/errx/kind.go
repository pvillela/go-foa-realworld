package errx

import (
	"github.com/pkg/errors"
)

type Kind struct {
	msg string
}

func KindOf(err error) *Kind {
	errx, ok := err.(*errxImpl)
	if !ok {
		return nil
	}
	return errx.kind
}

func NewKind(msg string) *Kind {
	return &Kind{msg}
}

func (s *Kind) Make(cause error, args ...interface{}) Errx {
	err := errxImpl{kind: s}
	err.args = args
	err.cause = cause
	err.tracer = errors.WithStack(dummyError{}).(stackTracer)
	return &err
}

func (s *Kind) Decorate(cause Errx, args ...interface{}) Errx {
	err := errxImpl{kind: s}
	err.args = args
	err.cause = cause
	return &err
}
