package errx

import (
	"fmt"
	"runtime/debug"
)

// Kind encapsulates an error messag string (possibly with argument placeholders) and a list
// of direct super-kinds it is to be considered to be related to.
type Kind struct {
	msg              string
	directSuperKinds []*Kind
}

// KindOf returns a pointer to the kind of an Errx or nil if err is not an Errx.
func KindOf(err error) *Kind {
	errx, ok := err.(*errxImpl)
	if !ok {
		return nil
	}
	return errx.kind
}

// NewKind instantiates a Kind.
func NewKind(msg string, superKinds ...*Kind) *Kind {
	return &Kind{msg, superKinds}
}

// Helper function to instantiate Errx instances from Kind pointers.
func (s *Kind) makeInternal(
	cause error,
	stackLinesToSuppress int,
	args ...any,
) *errxImpl {
	err := errxImpl{kind: s}
	err.args = args
	err.cause = cause
	err.stack = debug.Stack()
	err.stackLinesToSuppress = stackLinesToSuppress
	return &err
}

// Make instantiates an Errx from a Kind pointer, creating a stack trace at the point of
// instantiation.
func (s *Kind) Make(cause error, args ...any) Errx {
	err := s.makeInternal(cause, 3, args...)
	return err
}

// Decorate instantiates an Errx from a Kind pointer. The cause argument must be an Errx.
// The difference between Make and Decorate is that Make sets a new stack trace at the point
// of instantiation while Decorate effectively relies on the the stack trace provided by
// causes or its most recent cause which has a stack trace.
func (s *Kind) Decorate(cause Errx, args ...any) Errx {
	err := errxImpl{kind: s}
	err.args = args
	err.cause = cause
	return &err
}

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

// Helper function to add a Kind pointer to a set of Kind pointers.
func put(m map[*Kind]struct{}, kind *Kind) {
	m[kind] = struct{}{}
}

// SuperKinds returns the set of all super-kinds of the receiver, including the receiver itself.
func (s *Kind) SuperKinds() map[*Kind]struct{} {
	result := make(map[*Kind]struct{}, len(s.directSuperKinds))
	seen := make(map[*Kind]struct{}, len(s.directSuperKinds))
	workQueue := make([]*Kind, 1, len(s.directSuperKinds)+1)
	workQueue[0] = s
	for len(workQueue) > 0 {
		kind := workQueue[0]
		if _, ok := seen[kind]; !ok {
			put(seen, kind)
			put(result, kind)
			workQueue = append(workQueue[1:], kind.directSuperKinds...)
		} else {
			workQueue = workQueue[1:]
		}
	}
	return result
}

func (s *Kind) IsSubKindOf(kind *Kind) bool {
	superKinds := s.SuperKinds()
	_, ok := superKinds[kind]
	return ok
}
