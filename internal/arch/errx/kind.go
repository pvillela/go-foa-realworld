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
	if cause == nil {
		err.errorsCompanion = errors.New(err.msgWithArgs())
	} else if causex, ok := cause.(*errxImpl); ok {
		err.errorsCompanion = errors.Wrap(causex.errorsCompanion, err.msgWithArgs())
	} else {
		err.errorsCompanion = errors.Wrap(cause, err.msgWithArgs())
	}
	return &err
}

func (s Kind) Decorate(cause Errx, args ...interface{}) Errx {
	err := errxImpl{kind: s}
	err.args = args
	err.cause = cause
	if cause == nil {
		return s.Make(cause, args)
	} else if causex, ok := cause.(*errxImpl); ok {
		err.errorsCompanion = errors.WithMessage(causex.errorsCompanion, err.msgWithArgs())
	} else {
		// This can't happen.
		panic("errxImpl is the only valid implementation of interface Errx")
	}
	return &err
}
